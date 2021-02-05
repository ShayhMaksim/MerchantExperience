package processing

import (
	"database/sql"
	"strconv"
)

type rensponsibility struct {
	seller  seller
	product product
}

// Получение вьхи продавцов и товаров из БД
func getViewRensposibility(db *sql.DB) []rensponsibility {
	m_lrensposibility := []rensponsibility{}
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

//Добавление новых данных в БД
func addProducts(db *sql.DB, seller_id uint64, products []product) uint {
	var added uint = 0 //счетчик добавленных строк
	lenght := len(products)
	if lenght == 0 {
		return added
	}

	//мне кажется, что лучше составить один большой текстовый запрос
	for _, value := range products {
		productExec, err := db.Exec("insert into products (offer_id, name, price, quantity) values ($1, $2, $3, $4)",
			value.offer_id, value.name, value.price, value.quantity)
		checkErr(err)
		result, _ := productExec.RowsAffected()
		added += uint(result)
		sellerExec, err := db.Exec("insert into sellers (seller_id,offer_id) values ($1, $2)",
			seller_id, value.offer_id)
		checkErr(err)
		_, _ = sellerExec.RowsAffected()
	}
	return added
}

//удаление данных из БД
func deleteProducts(db *sql.DB, seller_id uint64, products []product) uint {
	var deleted uint = 0 // счетчик удаленных товаров
	lenght := len(products)
	if lenght == 0 {
		return deleted
	}

	for _, value := range products {
		productExec, err := db.Exec("delete from products where offer_id=$1",
			value.offer_id)
		checkErr(err)

		result, _ := productExec.RowsAffected()
		deleted += uint(result)
		sellerExec, err := db.Exec("delete from sellers where  seller_id=$1 and offer_id=$2",
			seller_id, value.offer_id)
		checkErr(err)
		_, _ = sellerExec.RowsAffected()
	}

	return deleted
}

//Обновление данных в БД
func updateProducts(db *sql.DB, products []product) uint {
	var updated uint = 0 // счетчик обновленных товаров
	lenght := len(products)
	if lenght == 0 {
		return updated
	}

	for _, value := range products {
		productExec, err := db.Exec("update products set name=$2, price=$3, quantity=$4 where offer_id=$1",
			value.offer_id, value.name, value.price, value.quantity)
		checkErr(err)
		result, _ := productExec.RowsAffected()
		updated += uint(result)
	}
	return updated
}

//Получение куска данных из БД
func localSelect(db *sql.DB, seller_id uint64, offer_id uint64, name string) []rensponsibility {
	m_lrensposibility := []rensponsibility{}

	var rows *sql.Rows
	var err error

	query := string("select * from responsibility where name like '" + name + "%'")

	seller_id_str := strconv.FormatUint(uint64(seller_id), 10)
	offer_id_str := strconv.FormatUint(uint64(offer_id), 10)

	if seller_id != 0 {
		query += " and seller_id=" + seller_id_str
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
