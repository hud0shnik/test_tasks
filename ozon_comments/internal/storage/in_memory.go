package storage

import (
	"comments/graph/model"
	"fmt"
	"sync"
	"time"
)

// InMemoryStorage - структура для хранения постов и комментариев
type InMemoryStorage struct {
	Post    *InMemoryPostStorage
	Comment *InMemoryCommentStorage
}

// InMemoryPostStorage - структура для хранения постов
type InMemoryPostStorage struct {
	mu     sync.RWMutex
	lastID int
	posts  []model.Post
}

// InMemoryCommentStorage - структура для хранения комментариев
type InMemoryCommentStorage struct {
	mu       sync.RWMutex
	lastID   int
	comments []model.Comment
}

// NewInMemoryPostStorage создаёт InMemoryPostStorage для хранения постов в памяти
// параметр cap - capacity слайса, в котором будут храниться посты
func NewInMemoryPostStorage(cap int) *InMemoryPostStorage {
	return &InMemoryPostStorage{
		lastID: 0,
		posts:  make([]model.Post, 0, cap),
	}
}

// NewInMemoryCommentStorage создаёт InMemoryCommentStorage для хранения комментариев в памяти
// параметр cap - capacity слайса, в котором будут храниться комментарии
func NewInMemoryCommentStorage(cap int) *InMemoryCommentStorage {
	return &InMemoryCommentStorage{
		lastID:   0,
		comments: make([]model.Comment, 0, cap),
	}
}

// AddPost - метод для добавляения поста в InMemoryPostStorage
func (s *InMemoryPostStorage) AddPost(post model.Post) (*model.Post, error) {

	// Блокируем мьютекс и дефером разлочиваем его
	s.mu.Lock()
	defer s.mu.Unlock()

	// Инкриментируем крайний айди поста
	s.lastID++

	// Задаём для нового поста айди и время создания
	post.ID = s.lastID
	post.CreatedAt = time.Now().GoString()

	// Добавляем новый пост
	s.posts = append(s.posts, post)

	return &post, nil
}

// GetPostById - метод поиска поста по айди
func (s *InMemoryPostStorage) GetPostById(id int) (model.Post, error) {

	// Блокируем мьютекс и дефером разлочиваем его
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Проверка на введённый айди
	if id > s.lastID || id <= 0 {
		return model.Post{}, fmt.Errorf("bad request")
	}

	return s.posts[id-1], nil
}

// GetAllPosts - метод получения всем постов.
// Параметры left и right - границы вывода (для пагинации)
func (s *InMemoryPostStorage) GetAllPosts(left, right int) ([]model.Post, error) {

	// Блокируем мьютекс и дефером разлочиваем его
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Проверка на выход айди из границ
	if left > s.lastID {
		return nil, nil
	}

	// Проверка на вывод без ограничений справа
	if left+right > s.lastID || right == -1 {
		return s.posts[left:], nil
	}

	// Проверка на некорректный ввод
	if left < 0 || right < 0 {
		return nil, fmt.Errorf("bad request")
	}

	return s.posts[left : left+right], nil
}

// AddComment - метод для добавляения комментария в InMemoryCommentStorage
func (s *InMemoryCommentStorage) AddComment(input model.CommentIntput) (*model.Comment, error) {

	// Блокируем мьютекс и дефером разлочиваем его
	s.mu.Lock()
	defer s.mu.Unlock()

	newComment := model.Comment{
		Post:    input.Post,
		Author:  input.Author,
		Content: input.Content,
		ReplyTo: input.ReplyTo,
	}

	// Инкриментируем крайний айди комментария
	s.lastID++

	// Задаём для нового коментария айди и время создания
	newComment.ID = s.lastID
	newComment.CreatedAt = time.Now().String()

	// Добавляем новый комментарий
	s.comments = append(s.comments, newComment)

	return &newComment, nil

}

// GetCommentsByPost - метод поиска комментариев у поста
func (s *InMemoryCommentStorage) GetCommentsByPost(postId, left, right int) ([]*model.Comment, error) {

	// Блокируем мьютекс и дефером разлочиваем его
	s.mu.RLock()
	defer s.mu.RUnlock()

	var res []*model.Comment

	// Проходим по всем комментариям и добавляем нужные в результат
	for _, comment := range s.comments {
		if comment.ReplyTo == nil && comment.Post == postId {
			com := comment
			res = append(res, &com)
		}
	}

	// Проверка на несуществующую страницу
	if left > len(res) {
		return nil, nil
	}

	// Проверка на вывод без ограничений справа
	if left+right > len(res) || right == -1 {
		return res[left:], nil
	}

	// Проверка на некорректный ввод
	if left < 0 || right < 0 {
		return nil, fmt.Errorf("bad request")
	}

	return res[left : left+right], nil
}

// GetRepliesOfComment - метод получения всех ответов на комментарий по id
func (s *InMemoryCommentStorage) GetRepliesOfComment(id int) ([]*model.Comment, error) {

	// Блокируем мьютекс и дефером разлочиваем его
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Проверка на выход айди из пределов
	if id > s.lastID {
		return nil, nil
	}

	var res []*model.Comment

	// Проходим по всем комментариям и добавляем нужные в результат
	for _, comment := range s.comments {
		if comment.ReplyTo != nil && *comment.ReplyTo == id {
			res = append(res, &comment)
		}
	}

	return res, nil
}
