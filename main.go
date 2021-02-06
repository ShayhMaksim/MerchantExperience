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

var database *sql.DB

func updateNewData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var inputData InputData
	_ = json.NewDecoder(r.Body).Decode(&inputData)

	fmt.Println(inputData.Selled_id)
	fmt.Println(inputData.ExcelFileName)

	xlsxData := processing.ReadDataFromXLSX(inputData.ExcelFileName)

	del := processing.DelegateRequest(database, inputData.Selled_id, xlsxData)

	json.NewEncoder(w).Encode(del)
}

func main() {

	database = processing.InitDatabase()
	r := mux.NewRouter()
	r.HandleFunc("/data", updateNewData).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", r))
}
