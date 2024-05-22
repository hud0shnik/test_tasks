package storage

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

// GetPostgresDB - функция подключения к Postgres
func GetPostgresDB() (*sqlx.DB, error) {

	// Подключение к БД
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PASSWORD")))
	if err != nil {
		return nil, fmt.Errorf("in sqlx.Open: %w", err)
	}

	// Проверка подключения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("in db.Ping: %w", err)
	}

	return db, nil
}
