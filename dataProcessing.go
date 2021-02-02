package main

import (
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
			// for col := 0; col < sheet.MaxCol; col++ {
			// 	text, err := sheet.Cell(row, col)
			// }

			// var offer_id uint
			// var name string
			// var price float32
			// var quantity int
			// var available bool

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
