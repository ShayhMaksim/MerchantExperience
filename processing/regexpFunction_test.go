package processing

import (
	"fmt"
	"testing"
)

type testOrderId struct {
	intput  string
	result  uint64
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
	result  uint64
	correct bool
}

type testAvailable struct {
	intput  string
	result  bool
	correct bool
}

var testsOrderId = []testOrderId{
	{"1234", 1234, true},
	{"1234A", 0, false},
	{"A12340", 0, false},
	{"-12340", 0, false},
	{"01234", 0, false},
}

// var testsName = []testName{
// 	{"1234", "", false},
// 	{"1234A", "", false},
// 	{"A12340", "A12340", true},
// 	{"-12340", "", false},
// 	{"01234", "", false},
// 	{"телевизор", "телевизор", true},
// 	{"телевизор A100", "телевизор A100", true},
// }

func TestsCorrectOfferId(t *testing.T) {
	fmt.Println("Тестирование IsCorrectOfferId ")
	for index, pair := range testsOrderId {
		result, correct := IsCorrectOfferId(pair.intput)
		if (result != pair.result) && (correct != pair.correct) {
			t.Error(
				"In the index", index,
				"For", pair.correct,
				"expected", pair.result, pair.correct,
				"got", result, correct,
			)
		}
	}
}

// func TestsCorrectName(t *testing.T) {
// 	fmt.Println("Тестирование IsCorrectName ")
// 	for index, pair := range testsName {
// 		result, correct := IsCorrectName(pair.intput)
// 		if (result != pair.result) && (correct != pair.correct) {
// 			t.Error(
// 				"In the index", index,
// 				"For", pair.correct,
// 				"expected", pair.result, pair.correct,
// 				"got", result, correct,
// 			)
// 		}
// 	}
// }
