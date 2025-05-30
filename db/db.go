package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Ошибка при подключении к базе:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("❌ База не отвечает:", err)
	}

	fmt.Println("✅ Успешное подключение к PostgreSQL на Render!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🚀 Go + PostgreSQL на Render работает!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Сервер запущен на порту " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
