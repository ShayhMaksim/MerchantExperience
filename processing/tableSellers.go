package processing

import "database/sql"

type Seller struct {
	Seller_id uint64 `json:"seller_id"`
	Offer_id  uint64 `json:"offer_id"`
}

// Получение всех продавцов из БД
func GetAllSellers(db *sql.DB) []Seller {
	Sellers := []Seller{}
	rows, err := db.Query("select * from Sellers")
	checkErr(err)
	for rows.Next() {

		var Seller_id uint64
		var Offer_id uint64

		err = rows.Scan(&Seller_id, &Offer_id)
		Sellers = append(Sellers,
			Seller{
				Seller_id: Seller_id,
				Offer_id:  Offer_id,
			})
	}
	return Sellers
}
