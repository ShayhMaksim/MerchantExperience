package main

import "fmt"

//"regexp"

func main() {
	db := initDatabase()

	products := localSelect(db, 4, 0, "теле")

	for _, value := range products {
		fmt.Println(value.product.offer_id, value.product.name, value.product.price, value.product.quantity, value.product.available)
	}

	// products := getAllProducts(db)

	// products := readDataFromXLSX(".\\Excel\\Exel.xlsx")
	// addProducts(db, 4, products)
	// //deleteProducts(db, 4, products)
	// updateProducts(db, products)
	// fmt.Println(products[0].price)

}
