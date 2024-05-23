package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"comments/graph/model"
	"comments/internal/storage"
	"context"
	"fmt"
)

var MAX_CONTENT_LENGTH = 2000

// CreateComment is the resolver for the CreateComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CommentIntput) (*model.Comment, error) {

	// Проверка вводимых параметров
	switch true {
	case len(input.Content) >= MAX_CONTENT_LENGTH:
		return nil, fmt.Errorf("too long content length")
	case len(input.Author) == 0:
		return nil, fmt.Errorf("no author")
	case input.Post < 1:
		return nil, fmt.Errorf("wrong post id")
	default:
	}

	if r.Posgres != nil {
		resutl, err := storage.SaveCommentToPostgres(r.Posgres, input)
		if err != nil {
			return nil, err
		}
		return &resutl, nil
	} else {
		result, err := r.InMemoryStorage.Comment.AddComment(input)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

}
