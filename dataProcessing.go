package main

import (
	"fmt"

	"github.com/tealeg/xlsx"
)

//почему-то Баг с отсутсвием row до сих пор не пофиксили -_-
func readDataFromXLSX(exelFileName string) {
	xlFile, err := xlsx.OpenFile(exelFileName)
	checkErr(err)
	for _, sheet := range xlFile.Sheets {
		fmt.Println(&sheet)
		for row := 0; row != sheet.MaxRow; row++ {
			for col := 0; col < sheet.MaxCol; col++ {
				text, err := sheet.Cell(row, col)
				checkErr(err)
				fmt.Println(text)
			}
		}
	}

}
