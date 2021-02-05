package processing

import (
	"database/sql"
	"fmt"

	"github.com/tealeg/xlsx"
)

//краткая статистика
// Added - количество созданных товаров
// Updated - количество обновленных товаров
// Deleted - количество удаленных товаров
// Wrong - количество неправильных строк
type Declaration struct {
	Added   uint
	Updated uint
	Deleted uint
	Wrong   uint
}

// структура, возвращающая результат с проверкой на корректность данных
type XlsxData struct {
	Product Product
	Correct bool
}

//проверка на корректность считанных данных
func IsCorrect(lOffer_id, lName, lPrice, lQuantity, lAvailable string) XlsxData {

	isCorrect := true
	prod := Product{}
	//проверка на положительность чисел
	prod.Offer_id, isCorrect = IsCorrectOfferId(lOffer_id)

	//проверка на отсутсвите в начале чисел (товар не должен начинаться с чисел?)
	prod.Name, isCorrect = IsCorrectName(lName)

	//проверка на отсутствие знаков и букв лишних в числе с плаваюещей точкой
	prod.Price, isCorrect = IsCorrectPrice(lPrice)

	//проверка на положительность чисел
	prod.Quantity, isCorrect = IsCorrectQuantity(lQuantity)

	//проверка на правильной записи типа bool
	prod.Available, isCorrect = IsCorrectAvailable(lAvailable)

	return XlsxData{prod, isCorrect}
}

//почему-то Баг с отсутсвием row до сих пор не пофиксили -_-
func ReadDataFromXLSX(exelFileName string) []XlsxData {
	xlsxDatas := []XlsxData{}

	xlFile, err := xlsx.OpenFile(exelFileName)
	checkErr(err)
	for _, sheet := range xlFile.Sheets {
		fmt.Println(&sheet)
		for row := 0; row != sheet.MaxRow; row++ {

			//На тот случай, если у нас данные сдвинуты в таблице по строке
			position := 0
			lOffer_id, _ := sheet.Cell(row, position+0)
			for lOffer_id == nil {
				position++
			}

			// lOffer_id, _ := sheet.Cell(row, position+0)
			lName, _ := sheet.Cell(row, position+1)
			lPrice, _ := sheet.Cell(row, position+2)
			lQuantity, _ := sheet.Cell(row, position+3)
			lAvailable, _ := sheet.Cell(row, position+4)

			xlsxData := IsCorrect(lOffer_id.Value, lName.Value, lPrice.Value, lQuantity.Value, lAvailable.Value)
			xlsxDatas = append(xlsxDatas, xlsxData)
		}
	}
	return xlsxDatas
}

/*
Общее обновление согласно поданным данным;
   declaration - статистика по данным:
Added - добавлено
Updated - обновлено
Deleted - удалено
Wrong - ошиблись
*/
func DelegateRequest(db *sql.DB, Seller_id uint64, Products []XlsxData) Declaration {
	addForProducts := []Product{}
	updateForProducts := []Product{}
	deleteForProducts := []Product{}
	//rensponsibilities := getViewRensposibility(db)
	var Wrong uint = 0

	for _, value := range Products {
		isUpdated := false //проверка на обновление данных

		//проверка на корректность данных
		if value.Correct == false {
			Wrong++
			continue
		}

		rensponsibilities := LocalSelect(db, Seller_id, value.Product.Offer_id, value.Product.Name)

		//обновление данных происходит в том случае, если указанный id продавца совпадает с id продавца из БД
		for _, Rensponsibility := range rensponsibilities {
			if Rensponsibility.Seller.Offer_id == value.Product.Offer_id && Seller_id == Rensponsibility.Seller.Seller_id {
				//обновление данных
				UpdatedProduct := Product{}
				if value.Product.Available == true {
					//просто занимаемся сложением данных
					UpdatedProduct.Quantity = Rensponsibility.Product.Quantity + value.Product.Quantity
					UpdatedProduct.Offer_id = value.Product.Offer_id
					UpdatedProduct.Name = value.Product.Name
					UpdatedProduct.Price = value.Product.Price
					UpdatedProduct.Available = true
					updateForProducts = append(updateForProducts, UpdatedProduct)
				}

				if value.Product.Available == false {
					UpdatedProduct.Quantity = Rensponsibility.Product.Quantity - value.Product.Quantity
					if UpdatedProduct.Quantity > 0 {
						//если товаров больше 0, то просто обновляем данные
						UpdatedProduct.Offer_id = value.Product.Offer_id
						UpdatedProduct.Name = value.Product.Name
						UpdatedProduct.Price = value.Product.Price
						UpdatedProduct.Available = value.Product.Available
						updateForProducts = append(updateForProducts, UpdatedProduct)
					} else if UpdatedProduct.Quantity == 0 {
						// если товаров не осталось, то нужно просто удалить из БД
						deleteForProducts = append(deleteForProducts, Rensponsibility.Product)
					} else if UpdatedProduct.Quantity < 0 {
						//если в excel указано, что у нас больше товаров идет на удаление, то тут какая-то ошибка
						Wrong++
					}
				}
				isUpdated = true
				//break
			}
		}

		if isUpdated == false {
			if value.Product.Available == true {
				addForProducts = append(addForProducts, value.Product)
			} else {
				Wrong++
			}
		}
	}

	Added := AddProducts(db, Seller_id, addForProducts)
	Updated := UpdateProducts(db, updateForProducts)
	Deleted := DeleteProducts(db, Seller_id, deleteForProducts)

	declaration := Declaration{
		Added:   Added,
		Updated: Updated,
		Deleted: Deleted,
		Wrong:   Wrong,
	}

	return declaration
}

//глубокая идея с подменой данных
func IsSimilar(dbProduct, excelProduct Product) bool {
	isSimilar := true
	if dbProduct.Name != excelProduct.Name {
		isSimilar = false
	}
	if dbProduct.Price != excelProduct.Price {
		isSimilar = false
	}
	return isSimilar
}
