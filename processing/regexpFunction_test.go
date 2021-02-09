package processing

import (
	"fmt"
	"testing"
)

type testOrderId struct {
	intput  string
	result  uint
	correct bool
}

type testName struct {
	intput  string
	result  string
	correct bool
}

type testPrice struct {
	intput  string
	result  float32
	correct bool
}

type testQuantity struct {
	intput  string
	result  uint
	correct bool
}

type testAvailable struct {
	intput  string
	result  bool
	correct bool
}

var testsOrderId = []testOrderId{
	{"1", 1, true},
	{"1234", 1234, true},
	{"1234A", 0, false},
	{"A12340", 0, false},
	{"-12340", 0, false},
	{"01234", 0, false},
}

var testsName = []testName{
	{"1234", "", false},
	{"1234A", "", false},
	{"A12340", "A12340", true},
	{"-12340", "", false},
	{"01234", "", false},
	{"телевизор", "телевизор", true},
	{"телевизор A100", "телевизор A100", true},
}

var testsPrice = []testPrice{
	{"1234", 1234, true},
	{"1234A", 0, false},
	{"A12340", 0, false},
	{"-12340", 0, false},
	{"01234", 1234, true},
	{"01234.43", 01234.43, true},
	{"01234,43", 01234.43, true}, //запятые тоже учитываем
	{"-01234,43", 0, false},
}

var testsQuantity = []testQuantity{
	{"1", 1, true},
	{"01234", 0, false},
	{"1234A", 0, false},
	{"A12340", 0, false},
	{"-12340", 0, false},
	{"01234.43", 0, false},
	{"01234,43", 0, false},
	{"-01234,43", 0, false},
	{"124", 124, true},
}

var testsAvailable = []testAvailable{
	{"01234", false, false},
	{"1234A", false, false},
	{"A12340", false, false},
	{"-12340", false, false},
	{"01234.43", false, false},
	{"01234,43", false, false},
	{"-01234,43", false, false},
	{"false", false, true},
	{"true", true, true},
	{"true1", false, false},
}

func TestIsCorrectOfferId(t *testing.T) {
	fmt.Println("Проверка на корректность поля offer_id")
	for index, pair := range testsOrderId {
		fmt.Print("Проверка {", pair.correct, ",", pair.intput, ",", pair.result, "}")
		result, correct := IsCorrectOfferId(pair.intput)
		if result != pair.result {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.result, pair.correct,
				"got", result, correct,
			)
		}
		if correct != pair.correct {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.correct,
				"got", correct,
			)
		}
		fmt.Println(" - Успешно пройдено")
	}
}

func TestIsCorrectName(t *testing.T) {
	fmt.Println("Проверка на корректность поля name")
	for index, pair := range testsName {
		fmt.Print("Проверка {", pair.correct, ",", pair.intput, ",", pair.result, "}")
		result, correct := IsCorrectName(pair.intput)
		if result != pair.result {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.result, pair.correct,
				"got", result, correct,
			)
		}
		if correct != pair.correct {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.correct,
				"got", correct,
			)
		}
		fmt.Println(" - Успешно пройдено")
	}
}

func TestIsCorrectPrice(t *testing.T) {
	fmt.Println("Проверка на корректность поля price")
	for index, pair := range testsPrice {
		fmt.Print("Проверка {", pair.correct, ",", pair.intput, ",", pair.result, "}")
		result, correct := IsCorrectPrice(pair.intput)
		if result != pair.result {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.result,
				"got", result,
			)
		}
		if correct != pair.correct {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.correct,
				"got", correct,
			)
		}
		fmt.Println(" - Успешно пройдено")
	}
}

func TestIsCorrectQuantity(t *testing.T) {
	fmt.Println("Проверка на корректность поля quantity")
	for index, pair := range testsQuantity {
		fmt.Print("Проверка {", pair.correct, ",", pair.intput, ",", pair.result, "}")
		result, correct := IsCorrectQuantity(pair.intput)
		if result != pair.result {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.result,
				"got", result,
			)
		}
		if correct != pair.correct {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.correct,
				"got", correct,
			)
		}
		fmt.Println(" - Успешно пройдено")
	}
}

func TestIsCorrectAvailable(t *testing.T) {
	fmt.Println("Проверка на корректность поля available")
	for index, pair := range testsAvailable {
		fmt.Print("Проверка {", pair.correct, ",", pair.intput, ",", pair.result, "}")
		result, correct := IsCorrectAvailable(pair.intput)
		if result != pair.result {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.result,
				"got", result,
			)
		}
		if correct != pair.correct {
			t.Error(
				"In the index", index,
				"For", pair.intput,
				"expected", pair.correct,
				"got", correct,
			)
		}
		fmt.Println(" - Успешно пройдено")
	}
}
