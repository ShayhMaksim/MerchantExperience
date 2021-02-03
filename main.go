package main

import (
	"fmt"
	"regexp"
	//"regexp"
)

func main() {
	//db := initDatabase()

	// products := getAllProducts(db)

	// for _, value := range products {
	// 	fmt.Println(value.offer_id, value.name, value.price, value.quantity, value.available)
	// }

	// products := readDataFromXLSX(".\\Excel\\Exel.xlsx")
	// addProducts(db, 4, products)
	// //deleteProducts(db, 4, products)
	// updateProducts(db, products)
	// fmt.Println(products[0].price)

	a := "Galaxy a"
	fmt.Println(a)
	//matched, err := regexp.MatchString("^[0-9]+$", a) для цифр
	//matched, err := regexp.MatchString("^[A-Za-z]+", a) для текста
	matched, err := regexp.MatchString("^Galaxy+", a)
	fmt.Println(matched)
	checkErr(err)
	//select * from responsibility where name LIKE '%Galaxy%'

}
