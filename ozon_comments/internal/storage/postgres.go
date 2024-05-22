package storage

import (
	"comments/graph/model"
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

// SavePostToPostgres сохраняет пост в Postgres. Поля id и created_at проставляются самим Postgres.
func SavePostToPostgres(db *sqlx.DB, post model.Post) error {

	// Вставка поста в БД
	_, err := db.Exec("INSERT INTO post (author, header, content, comments_allowed) VALUES ($1, $2, $3, $4)",
		post.Author, post.Header, post.Content, post.CommentsAllowed)
	if err != nil {
		return err
	}

	return nil

}
