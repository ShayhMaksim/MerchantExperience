package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"

	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sargaras/MerchantExperience/processing"
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

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) (*os.File, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return out, err
}

//загрузка данных
func updateNewData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var inputData InputData

	_ = json.NewDecoder(r.Body).Decode(&inputData)

	SELLER_ID := inputData.Selled_id
	excelFileReaded, _ := xlsx.OpenFile(inputData.ExcelFileName)

	structCh := make(chan struct{})
	del := processing.Declaration{}

	conveyor[staticKey] = AsynchDeclaration{Declaration: &del, ChStruct: &structCh, flag: false}
	key := UniqueKey{staticKey}
	json.NewEncoder(w).Encode(key)
	staticKey++

	go asynchAct(SELLER_ID, excelFileReaded, &del, structCh)

}

//асинхронное выполнение всех вычислений (по правилам подали ID продавца и его файл .xlxs)
func asynchAct(seller_id uint64, excelFile *xlsx.File, declaration *processing.Declaration, ch chan struct{}) {
	defer close(ch)
	xlsxData := processing.ReadDataFromXLSX(excelFile)
	*declaration = processing.DelegateRequest(database, seller_id, xlsxData)
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
	})
	r.HandleFunc("/data", updateNewData).Methods("POST")
	r.HandleFunc("/data/{id:[0-9]+}", getUpdatedData).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))
}
