package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error
	DB, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/newlabsgo?parseTime=true")
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal("БД недоступна:", err)
	}
	log.Println("✅ Подключено к MySQL!")
}
