package processing

import (
	"regexp"
	"strconv"
	"strings"
)

//проверка на положительность чисел
func IsCorrectOfferId(lOffer_id string) (uint64, bool) {
	var Offer_id uint64
	var errStr error
	var isCorrect bool = true
	matched, err := regexp.MatchString("^[1-9][0-9]+$", lOffer_id)
	checkErr(err)
	if matched == false {
		isCorrect = false
		Offer_id = 0
	} else {
		Offer_id, errStr = strconv.ParseUint(lOffer_id, 10, 64)
		checkErr(errStr)
	}
	return Offer_id, isCorrect
}

//проверка на отсутсвите в начале чисел (товар не должен начинаться с чисел?)
func IsCorrectName(lName string) (string, bool) {
	var Name string
	var isCorrect bool = true

	matched, err := regexp.MatchString("^[a-zA-ZА-Яа-я]+", lName)
	checkErr(err)
	if matched == false {
		isCorrect = false
		Name = ""
	} else {
		Name = string(lName)
	}
	return Name, isCorrect
}

//проверка на отсутствие знаков и букв лишних в числе с плаваюещей точкой
func IsCorrectPrice(lprice string) (float32, bool) {
	var price float32
	var isCorrect bool = true

	r := regexp.MustCompile("\\s+")
	replace := r.ReplaceAllString(lprice, "")
	lprice_value := strings.Split(string(replace), "р.")[0]

	matched, err := regexp.MatchString("^[0-9]*[.,]?[0-9]+$", lprice_value)
	checkErr(err)
	if matched == false {
		isCorrect = false
		price = 0.
	} else {
		f, _ := strconv.ParseFloat(lprice_value, 32)
		price = float32(f)
	}
	return price, isCorrect
}

//проверка на положительность чисел
func IsCorrectQuantity(lQuantity string) (uint64, bool) {
	var Quantity uint64
	var errStr error
	var isCorrect bool = true

	matched, err := regexp.MatchString("^[0-9]+$", lQuantity)
	checkErr(err)
	if matched == false {
		isCorrect = false
		Quantity = 0
	} else {
		Quantity, errStr = strconv.ParseUint(lQuantity, 10, 64)
		checkErr(errStr)
	}
	return Quantity, isCorrect
}

//проверка на правильной записи типа bool
func IsCorrectAvailable(lAvailable string) (bool, bool) {
	var Available bool
	var errStr error
	var isCorrect bool = true
	if (lAvailable != "true") && (lAvailable != "false") {
		isCorrect = false
		Available = false //??
	} else {
		Available, errStr = strconv.ParseBool(lAvailable)
		checkErr(errStr)
	}
	return Available, isCorrect
}
