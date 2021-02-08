package processing

import (
	"database/sql"
)

// Структура данных для товаров из БД
type Product struct {
	Offer_id  uint64  `json:"offer_id"`
	Name      string  `json:"name"`
	Price     float32 `json:"price"`
	Quantity  uint64  `json:"quantity"`
	Available bool    `json:"available"`
}

// Получение всех товаров из БД
func GetAllProducts(db *sql.DB) []Product {
	Products := []Product{}
	rows, err := db.Query("select * from Products")
	checkErr(err)
	for rows.Next() {

		var Offer_id uint64
		var Name string
		var Price float32
		var Quantity uint64
		var Available bool = true

		err = rows.Scan(&Offer_id, &Name, &Price, &Quantity)
		Products = append(Products,
			Product{
				Offer_id:  Offer_id,
				Name:      Name,
				Price:     Price,
				Quantity:  Quantity,
				Available: Available})
	}
	return Products

}
