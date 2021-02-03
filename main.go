package MerchantExperience

import "fmt"

func main() {
	db := initDatabase()

	// products := getAllProducts(db)

	// for _, value := range products {
	// 	fmt.Println(value.offer_id, value.name, value.price, value.quantity, value.available)
	// }

	products := readDataFromXLSX(".\\Excel\\Exel.xlsx")
	addProducts(db, 4, products)
	//deleteProducts(db, 4, products)
	updateProducts(db, products)
	fmt.Println(products[0].price)

}
