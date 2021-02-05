package processing

import (
	"database/sql"
	"fmt"

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
	prod.offer_id, isCorrect = isCorrectOfferId(loffer_id)

	//проверка на отсутсвите в начале чисел (товар не должен начинаться с чисел?)
	prod.name, isCorrect = isCorrectName(lname)

	//проверка на отсутствие знаков и букв лишних в числе с плаваюещей точкой
	prod.price, isCorrect = icCorrectPrice(lprice)

	//проверка на положительность чисел
	prod.quantity, isCorrect = isCorrectQuantity(lquantity)

	//проверка на правильной записи типа bool
	prod.available, isCorrect = isCorrectAvailable(lavailable)

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

/*
Общее обновление согласно поданным данным;
   declaration - статистика по данным:
added - добавлено
updated - обновлено
deleted - удалено
wrong - ошиблись
*/
func delegateRequest(db *sql.DB, seller_id uint64, products []xlsxData) declaration {
	addForProducts := []product{}
	updateForProducts := []product{}
	deleteForProducts := []product{}
	//rensponsibilities := getViewRensposibility(db)
	var wrong uint = 0

	for _, value := range products {
		isUpdated := false //проверка на обновление данных

		//проверка на корректность данных
		if value.correct == false {
			wrong++
			continue
		}

		rensponsibilities := localSelect(db, seller_id, value.product.offer_id, value.product.name)

		//обновление данных происходит в том случае, если указанный id продавца совпадает с id продавца из БД
		for _, rensponsibility := range rensponsibilities {
			if rensponsibility.seller.offer_id == value.product.offer_id && seller_id == rensponsibility.seller.seller_id {
				//обновление данных
				updatedProduct := product{}
				if value.product.available == true {
					//просто занимаемся сложением данных
					updatedProduct.quantity = rensponsibility.product.quantity + value.product.quantity
					updatedProduct.offer_id = value.product.offer_id
					updatedProduct.name = value.product.name
					updatedProduct.price = value.product.price
					updatedProduct.available = true
					updateForProducts = append(updateForProducts, updatedProduct)
				}

				if value.product.available == false {
					updatedProduct.quantity = rensponsibility.product.quantity - value.product.quantity
					if updatedProduct.quantity > 0 {
						//если товаров больше 0, то просто обновляем данные
						updatedProduct.offer_id = value.product.offer_id
						updatedProduct.name = value.product.name
						updatedProduct.price = value.product.price
						updatedProduct.available = value.product.available
						updateForProducts = append(updateForProducts, updatedProduct)
					} else if updatedProduct.quantity == 0 {
						// если товаров не осталось, то нужно просто удалить из БД
						deleteForProducts = append(deleteForProducts, rensponsibility.product)
					} else if updatedProduct.quantity < 0 {
						//если в excel указано, что у нас больше товаров идет на удаление, то тут какая-то ошибка
						wrong++
					}
				}
				isUpdated = true
				//break
			}
		}

		if isUpdated == false {
			if value.product.available == true {
				addForProducts = append(addForProducts, value.product)
			} else {
				wrong++
			}
		}
	}

	added := addProducts(db, seller_id, addForProducts)
	updated := updateProducts(db, updateForProducts)
	deleted := deleteProducts(db, seller_id, deleteForProducts)

	declaration := declaration{
		added:   added,
		updated: updated,
		deleted: deleted,
		wrong:   wrong,
	}

	return declaration
}

//глубокая идея с подменой данных
func isSimilar(dbProduct, excelProduct product) bool {
	isSimilar := true
	if dbProduct.name != excelProduct.name {
		isSimilar = false
	}
	if dbProduct.price != excelProduct.price {
		isSimilar = false
	}
	return isSimilar
}
