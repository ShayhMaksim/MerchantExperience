package requests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sargaras/MerchantExperience/processing"
	"github.com/tealeg/xlsx"
)

var Database *sql.DB
var Conveyor map[uint]AsynchDeclaration // конвейер для асинхронной обработки
var StaticKey uint = 0                  //уникальный ключ для конвейера

type InputData struct {
	Selled_id     uint   `json:"selled_id"`
	ExcelFileName string `json:"excelFileName"`
}

type UniqueKey struct {
	ID uint `json:"id"`
}

type AsynchDeclaration struct {
	Declaration *processing.Declaration
	ChStruct    *chan struct{}
	flag        bool //факт, что данные получены
}

type ErrorJson struct {
	Error string `json:"error"`
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) (*os.File, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return out, err
}

//загрузка данных
func UpdateNewData(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var inputData InputData

	_ = json.NewDecoder(r.Body).Decode(&inputData)

	SELLER_ID := inputData.Selled_id
	excelFileReaded, err := xlsx.OpenFile(inputData.ExcelFileName)
	if err != nil {
		jerror := ErrorJson{"Нет такого файла в директории!"}
		json.NewEncoder(w).Encode(jerror)
		fmt.Println("Нет такого файла в директории!")
		return
	}

	structCh := make(chan struct{})
	del := processing.Declaration{}

	Conveyor[StaticKey] = AsynchDeclaration{Declaration: &del, ChStruct: &structCh, flag: false}
	key := UniqueKey{StaticKey}
	json.NewEncoder(w).Encode(key)
	StaticKey++

	go asynchAct(SELLER_ID, excelFileReaded, &del, &structCh)

}

//асинхронное выполнение всех вычислений (по правилам подали ID продавца и его файл .xlxs)
func asynchAct(seller_id uint, excelFile *xlsx.File, declaration *processing.Declaration, ch *chan struct{}) {
	defer close(*ch)
	xlsxData := processing.ReadDataFromXLSX(excelFile)
	*declaration = processing.DelegateRequest(Database, seller_id, xlsxData)
}

//получение обновленных данных
func GetUpdatedData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 32)
	uintid := uint(id)
	if val, ok := Conveyor[uintid]; ok {
		//do something here
		if val.flag == false {
			<-*val.ChStruct
			Conveyor[uintid] = AsynchDeclaration{val.Declaration, nil, true}
		}
		del := *val.Declaration
		json.NewEncoder(w).Encode(del)
		//delete(conveyor, id) // удаляем элемент из конвейера
	}
}

type InputInfoData struct {
	Selled_id uint   `json:"selled_id"`
	Offer_id  uint   `json:"offer_id"`
	Name      string `json:"name"`
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputInfoData InputInfoData

	_ = json.NewDecoder(r.Body).Decode(&inputInfoData)

	seller_id := inputInfoData.Selled_id
	offer_id := inputInfoData.Offer_id
	name := inputInfoData.Name
	selected := processing.LocalSelect(Database, seller_id, offer_id, name)
	json.NewEncoder(w).Encode(selected)
}
