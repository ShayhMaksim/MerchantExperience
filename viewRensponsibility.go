package main

import "database/sql"

type rensponsibility struct {
	seller  seller
	product product
}

// Получение всех товаров из БД
func getViewRensposibility(db *sql.DB) []rensponsibility {
	m_lrensposibility := []rensponsibility{}
	rows, err := db.Query("select * from responsibility")
	checkErr(err)
	for rows.Next() {

		var seller_id uint
		var offer_id uint
		var name string
		var price float32
		var quantity int
		var available bool

		err = rows.Scan(&seller_id, &offer_id, &name, &price, &quantity, &available)
		m_lrensposibility = append(m_lrensposibility,
			rensponsibility{
				seller: seller{
					seller_id: seller_id,
					offer_id:  offer_id,
				},
				product: product{
					offer_id:  offer_id,
					name:      name,
					price:     price,
					quantity:  quantity,
					available: available,
				}})
	}
	return m_lrensposibility

}
