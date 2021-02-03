package MerchantExperience

import "database/sql"

type seller struct {
	seller_id uint
	offer_id  uint
}

// Получение всех продавцов из БД
func getAllSellers(db *sql.DB) []seller {
	sellers := []seller{}
	rows, err := db.Query("select * from sellers")
	checkErr(err)
	for rows.Next() {

		var seller_id uint
		var offer_id uint

		err = rows.Scan(&seller_id, &offer_id)
		sellers = append(sellers,
			seller{
				seller_id: seller_id,
				offer_id:  offer_id,
			})
	}
	return sellers
}
