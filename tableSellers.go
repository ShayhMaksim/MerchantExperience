package main

import "database/sql"

type seller struct {
	seller_id uint64
	offer_id  uint64
}

// Получение всех продавцов из БД
func getAllSellers(db *sql.DB) []seller {
	sellers := []seller{}
	rows, err := db.Query("select * from sellers")
	checkErr(err)
	for rows.Next() {

		var seller_id uint64
		var offer_id uint64

		err = rows.Scan(&seller_id, &offer_id)
		sellers = append(sellers,
			seller{
				seller_id: seller_id,
				offer_id:  offer_id,
			})
	}
	return sellers
}
