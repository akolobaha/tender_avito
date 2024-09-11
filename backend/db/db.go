package db

import (
	"database/sql"
	"log"
	"log/slog"
	"tenders/internal/config"
)

func GetConnection() *sql.DB {
	db, err := sql.Open("postgres", config.ConnString)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		slog.Error(err.Error())
	}

	return db
}
