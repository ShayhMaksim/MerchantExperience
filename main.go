package main

import (
	"fmt"

	"./processing"
)

//"regexp"

func main() {

	db := processing.InitDatabase()

	// Products := localSelect(db, 4, 0, "теле")

	// for _, value := range Products {
	// 	fmt.Println(value.Product.offer_id, value.Product.name, value.Product.price, value.Product.quantity, value.Product.available)
	// }

	// Products := getAllProducts(db)

	Products := processing.ReadDataFromXLSX(".\\Excel\\Update.xlsx")
	pro := []processing.Product{}
	for _, value := range Products {
		pro = append(pro, value.Product)
	}
	// addProducts(db, 4, Products)
	del := processing.DelegateRequest(db, 4, Products)
	fmt.Println(del.added, del.deleted, del.updated, del.wrong)
	//deleteProducts(db, 4, pro)
	// updateProducts(db, Products)
	fmt.Println(Products[0].Product.price)

}
