package processing

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

/*
	Набор констант для строки входа в БД
*/
const (
	DB_USER     = "postgres"
	DB_PASSWORD = "12345"
	DB_NAME     = "Base"
)

/*
	Подключение к Базе данных
*/
func initDatabase() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

// Функция проверки на ошибку
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
