package MerchantExperience

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

//почему-то Баг с отсутсвием row до сих пор не пофиксили -_-
func readDataFromXLSX(exelFileName string) []product {
	products := []product{}

	xlFile, err := xlsx.OpenFile(exelFileName)
	checkErr(err)
	for _, sheet := range xlFile.Sheets {
		fmt.Println(&sheet)
		for row := 0; row != sheet.MaxRow; row++ {

			//На тот случай, если у нас данные сдвинуты в таблице по строке
			position := 0
			loffer_id, _ := sheet.Cell(row, position+0)
			for loffer_id == nil {
				position++
			}

			// loffer_id, _ := sheet.Cell(row, position+0)
			lname, _ := sheet.Cell(row, position+1)
			lprice, _ := sheet.Cell(row, position+2)
			lquantity, _ := sheet.Cell(row, position+3)
			lavailable, _ := sheet.Cell(row, position+4)

			offer_id, _ := strconv.ParseUint(loffer_id.Value, 10, 32)
			name := string(lname.Value)
			price, _ := strconv.ParseFloat(strings.Split(lprice.Value, "р.")[0], 32)
			quantity, _ := strconv.ParseInt(lquantity.Value, 10, 32)
			available, _ := strconv.ParseBool(lavailable.Value)

			products = append(products, product{
				offer_id:  uint(offer_id),
				name:      name,
				price:     float32(price),
				quantity:  int(quantity),
				available: available})
		}
	}
	return products
}

//Добавление новых данных в БД
func addProducts(db *sql.DB, seller_id uint, products []product) {
	for _, value := range products {
		productExec, err := db.Exec("insert into products (offer_id, name, price, quantity) values ($1, $2, $3, $4)",
			value.offer_id, value.name, value.price, value.quantity)
		checkErr(err)
		fmt.Println(productExec.RowsAffected())
		sellerExec, err := db.Exec("insert into sellers (seller_id,offer_id) values ($1, $2)",
			seller_id, value.offer_id)
		checkErr(err)
		fmt.Println(sellerExec.RowsAffected())
	}
}

//удаление данных из БД
func deleteProducts(db *sql.DB, seller_id uint, products []product) {
	for _, value := range products {
		productExec, err := db.Exec("delete from products where offer_id=$1",
			value.offer_id)
		checkErr(err)
		fmt.Println(productExec.RowsAffected())
		sellerExec, err := db.Exec("delete from sellers where  seller_id=$1 and offer_id=$2",
			seller_id, value.offer_id)
		checkErr(err)
		fmt.Println(sellerExec.RowsAffected())
	}
}

//Обновление данных в БД
func updateProducts(db *sql.DB, products []product) {
	for _, value := range products {
		productExec, err := db.Exec("update products set name=$2, price=$3, quantity=$4 where offer_id=$1",
			value.offer_id, value.name, value.price, value.quantity)
		checkErr(err)
		fmt.Println(productExec.RowsAffected())
	}
}

//Общее обновление согласно поданным данным
func delegateRequest(db *sql.DB, seller_id uint, products []product) {
	addForProducts := []product{}
	updateForProducts := []product{}
	rensponsibilitys := getViewRensposibility(db)
	for _, value := range products {
		if value.offer_id == rensponsibilitys.offer_id {
			updateForProducts = append(updateForProducts, value)
		} else {
			addForProducts = append(addForProducts, value)
		}
	}
	updateProducts(db, updateForProducts)
	addProducts(db, addForProducts)
}
