package processing

import "database/sql"

type Seller struct {
	seller_id uint64
	offer_id  uint64
}

// Получение всех продавцов из БД
func GetAllSellers(db *sql.DB) []Seller {
	sellers := []Seller{}
	rows, err := db.Query("select * from Sellers")
	checkErr(err)
	for rows.Next() {

		var seller_id uint64
		var offer_id uint64

		err = rows.Scan(&seller_id, &offer_id)
		sellers = append(sellers,
			Seller{
				seller_id: seller_id,
				offer_id:  offer_id,
			})
	}
	return sellers
}
