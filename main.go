package main

import (
	"./processing"
	"fmt"
)

//"regexp"

func main() {
	db := initDatabase()

	// products := localSelect(db, 4, 0, "теле")

	// for _, value := range products {
	// 	fmt.Println(value.product.offer_id, value.product.name, value.product.price, value.product.quantity, value.product.available)
	// }

	// products := getAllProducts(db)

	products := readDataFromXLSX(".\\Excel\\Update.xlsx")
	pro := []product{}
	for _, value := range products {
		pro = append(pro, value.product)
	}
	// addProducts(db, 4, products)
	del := delegateRequest(db, 4, products)
	fmt.Println(del.added, del.deleted, del.updated, del.wrong)
	//deleteProducts(db, 4, pro)
	// updateProducts(db, products)
	fmt.Println(products[0].product.price)

}
