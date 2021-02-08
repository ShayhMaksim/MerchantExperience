package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sargaras/MerchantExperience/processing"
	"github.com/sargaras/MerchantExperience/requests"
)

func main() {

	requests.Database = processing.InitDatabase()
	requests.Conveyor = make(map[uint64]requests.AsynchDeclaration)

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
	})
	r.HandleFunc("/data", requests.UpdateNewData).Methods("POST")
	r.HandleFunc("/data/{id:[0-9]+}", requests.GetUpdatedData).Methods("GET")
	r.HandleFunc("/info", requests.GetData).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
}
