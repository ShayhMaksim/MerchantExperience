package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

//краткая статистика
// added - количество созданных товаров
// updated - количество обновленных товаров
// deleted - количество удаленных товаров
// wrong - количество неправильных строк
type declaration struct {
	added   uint
	updated uint
	deleted uint
	wrong   uint
}

// структура, возвращающая результат с проверкой на корректность данных
type xlsxData struct {
	product product
	correct bool
}

//проверка на корректность считанных данных
func isCorrect(loffer_id, lname, lprice, lquantity, lavailable string) xlsxData {

	isCorrect := true
	prod := product{}
	//проверка на положительность чисел
	matched, err := regexp.MatchString("^[0-9]+$", loffer_id)
	checkErr(err)
	if matched == false {
		isCorrect = false
		prod.offer_id = 0
	} else {
		offer_id, _ := strconv.ParseUint(loffer_id, 10, 32)
		prod.offer_id = uint(offer_id)
	}

	//проверка на отсутсвите в начале чисел (товар не должен начинаться с чисел?)
	matched, err = regexp.MatchString("^[A-Za-z]+", lname)
	checkErr(err)
	if matched == false {
		isCorrect = false
		prod.name = ""
	} else {
		prod.name = string(lname)
	}

	//проверка на отсутствие знаков и букв лишних в числе с плаваюещей точкой
	lprice_value := strings.Split(lprice, "р.")[0]
	matched, err = regexp.MatchString("^[0-9]*[.]?[0-9]+$", lprice_value)
	checkErr(err)
	if matched == false {
		isCorrect = false
		prod.price = 0
	} else {
		price, _ := strconv.ParseFloat(lprice_value, 32)
		prod.price = float32(price)
	}

	//проверка на положительность чисел
	matched, err = regexp.MatchString("^[0-9]+$", lquantity)
	checkErr(err)
	if matched == false {
		isCorrect = false
		prod.quantity = 0
	} else {
		quantity, _ := strconv.ParseInt(lquantity, 10, 32)
		prod.quantity = int(quantity)
	}

	//проверка на правильной записи типа bool
	if (lavailable != "true") || (lavailable != "false") {
		isCorrect = false
		prod.available = false //??
	} else {
		available, _ := strconv.ParseBool(lavailable)
		prod.available = available
	}
	return xlsxData{prod, isCorrect}
}

//почему-то Баг с отсутсвием row до сих пор не пофиксили -_-
func readDataFromXLSX(exelFileName string) []xlsxData {
	xlsxDatas := []xlsxData{}

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

			xlsxData := isCorrect(loffer_id.Value, lname.Value, lprice.Value, lquantity.Value, lavailable.Value)
			xlsxDatas = append(xlsxDatas, xlsxData)
		}
	}
	return xlsxDatas
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
// func delegateRequest(db *sql.DB, seller_id uint, products []product) {
// 	addForProducts := []product{}
// 	updateForProducts := []product{}
// 	deleteForProducts := []product{}
// 	rowMistake := 0
// 	rensponsibilities := getViewRensposibility(db)
// 	for _, value := range products {

// 		for _, valueR := range rensponsibilities {
// 			if valueR.product.offer_id == value.offer_id {

// 			}
// 		}

// 		updateForProducts = append(updateForProducts, value)
// 		addForProducts = append(addForProducts, value)
// 	}
// 	// updateProducts(db, updateForProducts)
// 	// addProducts(db, addForProducts)
// 	// deleteProducts(db, deleteForProducts)
// }
