package main

import (
	"fmt"
)

func main() {
	db := initDatabase()

	products := getAllData(db)

	for _, value := range products {
		fmt.Println(value.offer_id, value.name, value.price, value.quantity, value.available)
	}

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Fprintf(w, "Hello World!")
	// })
	// http.ListenAndServe(":8080", nil)
}
