package main

import (
	"github.com/tealeg/xlsx"
)

func readDataFromXLSX(exelFileName string) {
	xlFile, err := xlsx.OpenFile(exelFileName)
	checkErr(err)
	for _, sheet := range xlFile.Sheets {

	}

}
