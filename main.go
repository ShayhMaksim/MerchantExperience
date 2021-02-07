package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"./processing"
	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"
)

type InputData struct {
	Selled_id     uint64 `json:"selled_id"`
	ExcelFileName string `json:"excelFileName"`
}

type UniqueKey struct {
	ID uint64 `json:"id"`
}

type AsynchDeclaration struct {
	Declaration *processing.Declaration
	ChStruct    *chan struct{}
	flag        bool //факт, что данные получены
}

var database *sql.DB
var conveyor map[uint64]AsynchDeclaration // конвейер для асинхронной обработки
var staticKey uint64 = 0                  //уникальный ключ для конвейера

//загрузка данных
func updateNewData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	fmt.Println("TEMP DIR:", os.TempDir())

	src, hdr, err := r.FormFile("my-file")

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(os.TempDir(), hdr.Filename))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer dst.Close()

	io.Copy(dst, src)

	SELLER_ID, _ := strconv.ParseUint(r.FormValue("seller-id"), 10, 64)

	structCh := make(chan struct{})
	del := processing.Declaration{}

	conveyor[staticKey] = AsynchDeclaration{Declaration: &del, ChStruct: &structCh, flag: false}
	key := UniqueKey{staticKey}
	json.NewEncoder(w).Encode(key)
	staticKey++

	file, _ := xlsx.OpenFile(dst.Name())

	go asynchAct(SELLER_ID, file, &del, structCh)

}

//получение данных
func getUpdatedData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	if val, ok := conveyor[id]; ok {
		//do something here
		if val.flag == false {
			<-*val.ChStruct
			conveyor[id] = AsynchDeclaration{val.Declaration, nil, true}
		}
		del := *val.Declaration
		json.NewEncoder(w).Encode(del)
		//delete(conveyor, id) // удаляем элемент из конвейера
	}
}

func asynchAct(seller_id uint64, file *xlsx.File, declaration *processing.Declaration, ch chan struct{}) {
	defer close(ch)
	// fmt.Println(inputData.Selled_id)
	// fmt.Println(inputData.ExcelFileName)
	xlsxData := processing.ReadDataFromXLSX(file)
	*declaration = processing.DelegateRequest(database, seller_id, xlsxData)
}

func main() {

	database = processing.InitDatabase()
	conveyor = make(map[uint64]AsynchDeclaration)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		http.ServeFile(w, r, ".\\user.html")
	})
	r.HandleFunc("/data", updateNewData).Methods("POST")
	r.HandleFunc("/data/{id}", getUpdatedData).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r))
}
