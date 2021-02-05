package processing

import (
	"database/sql"
)

// Структура данных для товаров из БД
type product struct {
	offer_id  uint64
	name      string
	price     float32
	quantity  uint64
	available bool
}

// Получение всех товаров из БД
func getAllProducts(db *sql.DB) []product {
	products := []product{}
	rows, err := db.Query("select * from products")
	checkErr(err)
	for rows.Next() {

		var offer_id uint64
		var name string
		var price float32
		var quantity uint64
		var available bool = true

		err = rows.Scan(&offer_id, &name, &price, &quantity)
		products = append(products,
			product{
				offer_id:  offer_id,
				name:      name,
				price:     price,
				quantity:  quantity,
				available: available})
	}
	return products

}
