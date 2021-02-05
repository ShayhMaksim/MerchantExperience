package processing

import (
	"database/sql"
	"strconv"
)

type Rensponsibility struct {
	seller  Seller
	product Product
}

// Получение вьхи продавцов и товаров из БД
func GetViewRensposibility(db *sql.DB) []Rensponsibility {
	m_lrensposibility := []Rensponsibility{}
	rows, err := db.Query("select * from responsibility")
	checkErr(err)
	for rows.Next() {

		var seller_id uint64
		var offer_id uint64
		var name string
		var price float32
		var quantity uint64
		var available bool = true

		err = rows.Scan(&seller_id, &offer_id, &name, &price, &quantity)
		m_lrensposibility = append(m_lrensposibility,
			Rensponsibility{
				seller: Seller{
					seller_id: seller_id,
					offer_id:  offer_id,
				},
				product: Product{
					offer_id:  offer_id,
					name:      name,
					price:     price,
					quantity:  quantity,
					available: available,
				}})
	}
	return m_lrensposibility

}

//Добавление новых данных в БД
func AddProducts(db *sql.DB, seller_id uint64, products []Product) uint {
	var added uint = 0 //счетчик добавленных строк
	lenght := len(products)
	if lenght == 0 {
		return added
	}

	//мне кажется, что лучше составить один большой текстовый запрос
	for _, value := range products {
		ProductExec, err := db.Exec("insert into Products (offer_id, name, price, quantity) values ($1, $2, $3, $4)",
			value.offer_id, value.name, value.price, value.quantity)
		checkErr(err)
		result, _ := ProductExec.RowsAffected()
		added += uint(result)
		SellerExec, err := db.Exec("insert into Sellers (Seller_id,offer_id) values ($1, $2)",
			seller_id, value.offer_id)
		checkErr(err)
		_, _ = SellerExec.RowsAffected()
	}
	return added
}

//удаление данных из БД
func DeleteProducts(db *sql.DB, seller_id uint64, products []Product) uint {
	var deleted uint = 0 // счетчик удаленных товаров
	lenght := len(products)
	if lenght == 0 {
		return deleted
	}

	for _, value := range products {
		ProductExec, err := db.Exec("delete from Products where offer_id=$1",
			value.offer_id)
		checkErr(err)

		result, _ := ProductExec.RowsAffected()
		deleted += uint(result)
		SellerExec, err := db.Exec("delete from Sellers where  Seller_id=$1 and offer_id=$2",
			seller_id, value.offer_id)
		checkErr(err)
		_, _ = SellerExec.RowsAffected()
	}

	return deleted
}

//Обновление данных в БД
func UpdateProducts(db *sql.DB, products []Product) uint {
	var updated uint = 0 // счетчик обновленных товаров
	lenght := len(products)
	if lenght == 0 {
		return updated
	}

	for _, value := range products {
		ProductExec, err := db.Exec("update Products set name=$2, price=$3, quantity=$4 where offer_id=$1",
			value.offer_id, value.name, value.price, value.quantity)
		checkErr(err)
		result, _ := ProductExec.RowsAffected()
		updated += uint(result)
	}
	return updated
}

//Получение куска данных из БД
func LocalSelect(db *sql.DB, seller_id uint64, offer_id uint64, name string) []Rensponsibility {
	m_lrensposibility := []Rensponsibility{}

	var rows *sql.Rows
	var err error

	query := string("select * from responsibility where name like '" + name + "%'")

	Seller_id_str := strconv.FormatUint(uint64(seller_id), 10)
	offer_id_str := strconv.FormatUint(uint64(offer_id), 10)

	if seller_id != 0 {
		query += " and Seller_id=" + Seller_id_str
	}
	if offer_id != 0 {
		query += " and offer_id=" + offer_id_str
	}

	rows, err = db.Query(query)

	checkErr(err)
	for rows.Next() {

		var seller_id uint64
		var offer_id uint64
		var name string
		var price float32
		var quantity uint64
		var available bool = true

		err = rows.Scan(&seller_id, &offer_id, &name, &price, &quantity)
		m_lrensposibility = append(m_lrensposibility,
			Rensponsibility{
				seller: Seller{
					seller_id: seller_id,
					offer_id:  offer_id,
				},
				product: Product{
					offer_id:  offer_id,
					name:      name,
					price:     price,
					quantity:  quantity,
					available: available,
				}})
	}
	return m_lrensposibility
}
