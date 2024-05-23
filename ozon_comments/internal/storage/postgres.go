package storage

import (
	"comments/graph/model"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
func SavePostToPostgres(db *sqlx.DB, post model.Post) (*model.Post, error) {

	// Вставка поста в БД и запись айди и времени создания
	res := db.QueryRow("INSERT INTO post (author, header, content, comments_allowed) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		post.Author, post.Header, post.Content, post.CommentsAllowed)
	err := res.Scan(&post.ID, &post.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &post, nil

}

// GetPostFromPostgresById ищет пост по айди в Postgres
func GetPostFromPostgresById(db *sqlx.DB, id *int) (model.Post, error) {

	var selected []model.Post

	// Поиск поста по айди
	err := db.Select(&selected, fmt.Sprintf("SELECT id, author, header, content, comments_allowed AS commentsAllowed, created_at AS createdAt FROM post WHERE id=%d LIMIT 1", *id))
	if err != nil {
		return model.Post{}, err
	}

	// Проверка на пустой результат поиска
	if len(selected) == 0 {
		return model.Post{}, fmt.Errorf("post not found")
	}

	return selected[0], nil
}

// GetAllPostsFromPostgres - функция получения всех постов.
// Параметры left и right определяют какую часть резльтата необходимо получить (пагинация)
func GetAllPostsFromPostgres(db *sqlx.DB, left, right int) ([]model.Post, error) {

	// Конструирование запроса к БД
	query := "SELECT id, author, header, content, comments_allowed AS commentsAllowed, created_at AS createdAt FROM post ORDER BY created_at OFFSET $1"
	args := []interface{}{left}
	if right > 0 {
		query += " LIMIT $2"
		args = append(args, right)
	}

	var posts []model.Post

	// Исполнение запроса и запись результата
	if err := db.Select(&posts, query, args...); err != nil {
		return nil, err
	}

	return posts, nil
}

// SaveCommentToPostgres - функция сохранения комментария в Postgres
func SaveCommentToPostgres(db *sqlx.DB, input model.CommentIntput) (model.Comment, error) {

	// Новый комментарий
	newComent := model.Comment{
		Post:    input.Post,
		ReplyTo: input.ReplyTo,
		Author:  input.Author,
		Content: input.Content,
	}

	// Формирование и исполнение запроса к Postgres
	row := db.QueryRow("INSERT INTO comment (post, author, content, reply_to) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		input.Post, input.Author, input.Content, input.ReplyTo)
	err := row.Scan(&newComent.ID, &newComent.CreatedAt)
	if err != nil {
		return model.Comment{}, err
	}

	return newComent, nil
}

// GetCommentsByPost - функция получения всех комментариев к посту.
// Параметры left и right определяют какую часть резльтата необходимо получить (пагинация)
func GetCommentsByPost(db *sqlx.DB, postId, left, right int) ([]*model.Comment, error) {

	// Формирование запроса к БД
	query := "SELECT id, post, author, content, created_at AS createdAt, reply_to AS replyTo FROM comment WHERE post = $1 AND reply_to IS NULL ORDER BY created_at OFFSET $2"
	args := []interface{}{postId, left}

	// Проверка на наличие ограничения справа
	if right >= 0 {
		query += " LIMIT $3"
		args = append(args, right)
	}

	// Выполнение запроса и запись результата
	var comments []*model.Comment
	if err := db.Select(&comments, query, args...); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetRepliesOfComment - функция получения всех ответов на комментарий по айди
func GetRepliesOfComment(db *sqlx.DB, id int) ([]*model.Comment, error) {

	var result []*model.Comment

	// Исполнение запроса и запись результата
	err := db.Select(&result, "SELECT id, post, author, content, created_at AS createdAt, reply_to AS replyTo FROM comment WHERE reply_to = $1", id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
