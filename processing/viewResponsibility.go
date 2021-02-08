package processing

import (
	"database/sql"
	"strconv"
)

type Rensponsibility struct {
	Seller  Seller
	Product Product
}

// Получение вьхи продавцов и товаров из БД
func GetViewRensposibility(db *sql.DB) []Rensponsibility {
	m_lrensposibility := []Rensponsibility{}
	rows, err := db.Query("select * from responsibility")
	checkErr(err)
	for rows.Next() {

		var Seller_id uint64
		var Offer_id uint64
		var Name string
		var Price float32
		var Quantity uint64
		var Available bool = true

		err = rows.Scan(&Seller_id, &Offer_id, &Name, &Price, &Quantity)
		m_lrensposibility = append(m_lrensposibility,
			Rensponsibility{
				Seller: Seller{
					Seller_id: Seller_id,
					Offer_id:  Offer_id,
				},
				Product: Product{
					Offer_id:  Offer_id,
					Name:      Name,
					Price:     Price,
					Quantity:  Quantity,
					Available: Available,
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
		ProductExec, err := db.Exec("insert into products (Offer_id, Name, Price, Quantity) values ($1, $2, $3, $4)",
			value.Offer_id, value.Name, value.Price, value.Quantity)
		checkErr(err)
		result, _ := ProductExec.RowsAffected()
		added += uint(result)
		SellerExec, err := db.Exec("insert into sellers (Seller_id,Offer_id) values ($1, $2)",
			seller_id, value.Offer_id)
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
		ProductExec, err := db.Exec("delete from products where Offer_id=$1",
			value.Offer_id)
		checkErr(err)

		result, _ := ProductExec.RowsAffected()
		deleted += uint(result)
		SellerExec, err := db.Exec("delete from sellers where  Seller_id=$1 and Offer_id=$2",
			seller_id, value.Offer_id)
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
		ProductExec, err := db.Exec("update products set Name=$2, Price=$3, Quantity=$4 where Offer_id=$1",
			value.Offer_id, value.Name, value.Price, value.Quantity)
		checkErr(err)
		result, _ := ProductExec.RowsAffected()
		updated += uint(result)
	}
	return updated
}

//Получение куска данных из БД
func LocalSelect(db *sql.DB, seller_id uint64, offer_id uint64, Name string) []Rensponsibility {
	m_lrensposibility := []Rensponsibility{}

	var rows *sql.Rows
	var err error

	query := string("select * from responsibility where Name like '" + Name + "%'")

	Seller_id_str := strconv.FormatUint(uint64(seller_id), 10)
	Offer_id_str := strconv.FormatUint(uint64(offer_id), 10)

	if seller_id != 0 {
		query += " and Seller_id=" + Seller_id_str
	}

	if offer_id != 0 {
		query += " and Offer_id=" + Offer_id_str
	}

	rows, err = db.Query(query)

	checkErr(err)
	for rows.Next() {

		var Seller_id uint64
		var Offer_id uint64
		var Name string
		var Price float32
		var Quantity uint64
		var Available bool = true

		err = rows.Scan(&Seller_id, &Offer_id, &Name, &Price, &Quantity)
		m_lrensposibility = append(m_lrensposibility,
			Rensponsibility{
				Seller: Seller{
					Seller_id: Seller_id,
					Offer_id:  Offer_id,
				},
				Product: Product{
					Offer_id:  Offer_id,
					Name:      Name,
					Price:     Price,
					Quantity:  Quantity,
					Available: Available,
				}})
	}
	return m_lrensposibility
}
