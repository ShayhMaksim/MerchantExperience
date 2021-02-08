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
	var excelFile *os.File
	var SELLER_ID uint64

	SELLER_ID, excelFile = readNewDataFromForm(w, r)

	structCh := make(chan struct{})
	del := processing.Declaration{}

	conveyor[staticKey] = AsynchDeclaration{Declaration: &del, ChStruct: &structCh, flag: false}
	key := UniqueKey{staticKey}
	json.NewEncoder(w).Encode(key)
	staticKey++

	go asynchAct(SELLER_ID, excelFile, &del, structCh)

}

//асинхронное выполнение всех вычислений (по правилам подали ID продавца и его файл .xlxs)
func asynchAct(seller_id uint64, file *os.File, declaration *processing.Declaration, ch chan struct{}) {
	defer close(ch)
	excelFile, _ := xlsx.OpenFile(file.Name())
	xlsxData := processing.ReadDataFromXLSX(excelFile)
	*declaration = processing.DelegateRequest(database, seller_id, xlsxData)
}

//чтение данных по форме (что по идее не должно быть???)
func readNewDataFromForm(w http.ResponseWriter, r *http.Request) (uint64, *os.File) {
	fmt.Println("TEMP DIR:", os.TempDir())

	src, hdr, err := r.FormFile("my-file")

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	defer src.Close()

	ExcelFile, err := os.Create(filepath.Join(os.TempDir(), hdr.Filename))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	defer ExcelFile.Close()

	io.Copy(ExcelFile, src)

	SELLER_ID, _ := strconv.ParseUint(r.FormValue("seller-id"), 10, 64)
	return SELLER_ID, ExcelFile
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
	log.Fatal(http.ListenAndServe(":8080", r))
}
