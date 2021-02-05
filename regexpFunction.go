package main

import (
	"regexp"
	"strconv"
	"strings"
)

//проверка на положительность чисел
func isCorrectOfferId(loffer_id string) (uint64, bool) {
	var offer_id uint64
	var errStr error
	var isCorrect bool = true
	matched, err := regexp.MatchString("^[0-9]+$", loffer_id)
	checkErr(err)
	if matched == false {
		isCorrect = false
		offer_id = 0
	} else {
		offer_id, errStr = strconv.ParseUint(loffer_id, 10, 64)
		checkErr(errStr)
		offer_id = uint(offer_id)
	}
	return offer_id, isCorrect
}

//проверка на отсутсвите в начале чисел (товар не должен начинаться с чисел?)
func isCorrectName(lname string) (string, bool) {
	var name string
	var isCorrect bool = true

	matched, err := regexp.MatchString("^[a-zA-ZА-Яа-я]+", lname)
	checkErr(err)
	if matched == false {
		isCorrect = false
		name = ""
	} else {
		name = string(lname)
	}
	return name, isCorrect
}

//проверка на отсутствие знаков и букв лишних в числе с плаваюещей точкой
func icCorrectPrice(lprice string) (float32, bool) {
	var price float32
	var errStr error
	var isCorrect bool = true

	r := regexp.MustCompile("\\s+")
	replace := r.ReplaceAllString(lprice, "")
	lprice_value := strings.Split(string(replace), "р.")[0]

	matched, err := regexp.MatchString("^[0-9]*[.,]?[0-9]+$", lprice_value)
	checkErr(err)
	if matched == false {
		isCorrect = false
		price = 0
	} else {
		price, errStr = strconv.ParseFloat(lprice_value, 32)
		checkErr(errStr)
		price = float32(price)
	}
	return price, isCorrect
}

//проверка на положительность чисел
func isCorrectQuantity(lquantity string) (uint64, bool) {
	var quantity uint64
	var errStr error
	var isCorrect bool = true

	matched, err := regexp.MatchString("^[0-9]+$", lquantity)
	checkErr(err)
	if matched == false {
		isCorrect = false
		quantity = 0
	} else {
		quantity, _ := strconv.ParseUint(lquantity, 10, 64)
		checkErr(errStr)
		quantity = quantity
	}
	return quantity, isCorrect
}

//проверка на правильной записи типа bool
func isCorrectAvailable(lavailable string) (bool, bool) {
	var available bool
	var errStr error
	var isCorrect bool = true
	if (lavailable != "true") && (lavailable != "false") {
		isCorrect = false
		available = false //??
	} else {
		available, _ := strconv.ParseBool(lavailable)
		checkErr(errStr)
		available = available
	}
	return available, isCorrect
}
