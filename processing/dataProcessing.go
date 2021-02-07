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
	Added   uint `json:"added"`
	Updated uint `json:"updated"`
	Deleted uint `json:"deleted"`
	Wrong   uint `json:"wrong"`
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
func ReadDataFromXLSX(xlFile *xlsx.File) []XlsxData {
	xlsxDatas := []XlsxData{}

	// xlFile, err := xlsx.OpenFile(exelFileName)
	// checkErr(err)
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
func DelegateRequest(db *sql.DB, seller_id uint64, xlsxData []XlsxData) Declaration {

	addForProducts := []Product{}
	updateForProducts := []Product{}
	deleteForProducts := []Product{}
	//rensponsibilities := getViewRensposibility(db)
	var wrong uint = 0

	for _, value := range xlsxData {
		isUpdated := false //проверка на обновление данных
		isCorrectSellerID := true
		//проверка на корректность данных
		if value.Correct == false {
			wrong++
			continue
		}

		rensponsibilities := LocalSelect(db, 0, value.Product.Offer_id, value.Product.Name)

		//обновление данных происходит в том случае, если указанный id продавца совпадает с id продавца из БД
		for _, rensponsibility := range rensponsibilities {
			//если айдишники совпадают по товару, а продавцы разные, то это ошибка
			if rensponsibility.Seller.Offer_id == value.Product.Offer_id {
				if rensponsibility.Seller.Seller_id != seller_id {
					wrong++
					isCorrectSellerID = false
					break //выход из цикла, потому что id товаров совпали
				}
			}

			if rensponsibility.Seller.Offer_id == value.Product.Offer_id && seller_id == rensponsibility.Seller.Seller_id {
				//обновление данных
				UpdatedProduct := Product{}
				if value.Product.Available == true {
					//просто занимаемся сложением данных
					UpdatedProduct.Quantity = rensponsibility.Product.Quantity + value.Product.Quantity
					UpdatedProduct.Offer_id = value.Product.Offer_id
					UpdatedProduct.Name = value.Product.Name
					UpdatedProduct.Price = value.Product.Price
					UpdatedProduct.Available = true
					updateForProducts = append(updateForProducts, UpdatedProduct)
				}

				if value.Product.Available == false {
					UpdatedProduct.Quantity = rensponsibility.Product.Quantity - value.Product.Quantity
					if UpdatedProduct.Quantity > 0 {
						//если товаров больше 0, то просто обновляем данные
						UpdatedProduct.Offer_id = value.Product.Offer_id
						UpdatedProduct.Name = value.Product.Name
						UpdatedProduct.Price = value.Product.Price
						UpdatedProduct.Available = value.Product.Available
						updateForProducts = append(updateForProducts, UpdatedProduct)
					} else if UpdatedProduct.Quantity == 0 {
						// если товаров не осталось, то нужно просто удалить из БД
						deleteForProducts = append(deleteForProducts, rensponsibility.Product)
					} else if UpdatedProduct.Quantity < 0 {
						//если в excel указано, что у нас больше товаров идет на удаление, то тут какая-то ошибка
						wrong++
					}
				}
				isUpdated = true
				break //выход из цикла, потому что id товаров совпали
			}
		}

		if (isUpdated == false) && (isCorrectSellerID == true) {
			if value.Product.Available == true {
				addForProducts = append(addForProducts, value.Product)
			} else {
				wrong++
			}
		}
	}

	added := AddProducts(db, seller_id, addForProducts)
	updated := UpdateProducts(db, updateForProducts)
	deleted := DeleteProducts(db, seller_id, deleteForProducts)

	declaration := Declaration{
		Added:   added,
		Updated: updated,
		Deleted: deleted,
		Wrong:   wrong,
	}
	return declaration
}
