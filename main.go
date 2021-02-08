package main

import (
	"database/sql"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sargaras/MerchantExperience/processing"
)



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
