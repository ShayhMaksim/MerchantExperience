package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	ChStruct    *chan struct{}
}

var database *sql.DB
var conveyor map[uint64]AsynchDeclaration // конвейер для асинхронной обработки
var staticKey uint64 = 0                  //уникальный ключ для конвейера

//загрузка данных
func updateNewData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputData InputData
	_ = json.NewDecoder(r.Body).Decode(&inputData)

	structCh := make(chan struct{})
	del := processing.Declaration{}

	conveyor[staticKey] = AsynchDeclaration{Declaration: &del, ChStruct: &structCh}
	key := UniqueKey{staticKey}
	json.NewEncoder(w).Encode(key)
	staticKey++

	go asynchAct(inputData, &del, structCh)

}

//получение данных
func getUpdatedData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 64)
	if val, ok := conveyor[id]; ok {
		//do something here
		<-*val.ChStruct
		del := *val.Declaration
		json.NewEncoder(w).Encode(del)
		//delete(conveyor, id) // удаляем элемент из конвейера
	}
}

func asynchAct(inputData InputData, declaration *processing.Declaration, ch chan struct{}) {
	defer close(ch)
	fmt.Println(inputData.Selled_id)
	fmt.Println(inputData.ExcelFileName)
	xlsxData := processing.ReadDataFromXLSX(inputData.ExcelFileName)
	*declaration = processing.DelegateRequest(database, inputData.Selled_id, xlsxData)
}

func main() {

	database = processing.InitDatabase()
	conveyor = make(map[uint64]AsynchDeclaration)

	r := mux.NewRouter()
	r.HandleFunc("/data", updateNewData).Methods("POST")
	r.HandleFunc("/data/{id}", getUpdatedData).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r))
}
