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

// GetPostFromPostgresById ищет пост по айди в Postgres
func GetPostFromPostgresById(db *sqlx.DB, id *int) (model.Post, error) {

	var selected []model.Post

	// Поиск поста по айди
	err := db.Select(&selected, fmt.Sprintf("SELECT id, author, header, content, comments_allowed AS commentsAllowed, created_at AS createdAt FROM post WHERE id=%d", *id))
	// Todo: добавить LIMIT 1
	if err != nil {
		return model.Post{}, err
	}

	if len(selected) == 0 {
		return model.Post{}, fmt.Errorf("post not found")
	}

	return selected[0], nil
}
