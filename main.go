package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"./processing"
	"github.com/gorilla/mux"
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
	chStruct    *chan struct{}
}

var database *sql.DB
var conveyor map[uint64]AsynchDeclaration // конвейер для асинхронной обработки
var staticKey uint64 = 0                  //уникальный ключ для конвейера

func updateNewData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputData InputData
	_ = json.NewDecoder(r.Body).Decode(&inputData)

	structCh := make(chan struct{})
	del := processing.Declaration{}

	conveyor[staticKey] = AsynchDeclaration{&del, &structCh}
	staticKey++

	fmt.Println(inputData.Selled_id)
	fmt.Println(inputData.ExcelFileName)

	go asynchAct(inputData, &del, structCh)
	<-structCh
	json.NewEncoder(w).Encode(del)
}

func getUpdatedData(w http.ResponseWriter, r *http.Request)
{
	w.Header().Set("Content-Type", "application/json")
	var inputData InputData
	_ = json.NewDecoder(r.Body).Decode(&inputData)
}

func asynchAct(inputData InputData, declaration *processing.Declaration, ch chan struct{}) {
	defer close(ch)
	xlsxData := processing.ReadDataFromXLSX(inputData.ExcelFileName)
	*declaration = processing.DelegateRequest(database, inputData.Selled_id, xlsxData)
}

func main() {

	database = processing.InitDatabase()
	conveyor = make(map[uint64]AsynchDeclaration)

	r := mux.NewRouter()
	r.HandleFunc("/data", updateNewData).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", r))
}
