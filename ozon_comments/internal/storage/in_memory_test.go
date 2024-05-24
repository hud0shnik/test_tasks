package storage

import (
	"comments/graph/model"
	"testing"
)

func isPostsEqual(first, second model.Post) bool {

	if first.Author == second.Author && first.Header == second.Header &&
		first.Content == second.Content && first.CommentsAllowed == second.CommentsAllowed {
		return true
	}
	return false
}

func TestAddPost(t *testing.T) {

	normalPost := model.Post{
		Author:          "normalAuthor",
		Header:          "Normal Header",
		Content:         "Normal content",
		CommentsAllowed: true,
	}

	tests := []struct {
		name    string
		input   model.Post
		want    *model.Post
		wantErr bool
	}{
		{
			name:    "Normal post check",
			input:   normalPost,
			want:    &normalPost,
			wantErr: false,
		},
		{
			name:    "Empty post check",
			input:   model.Post{},
			want:    &model.Post{},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := InMemoryStorage{
				Post:    NewInMemoryPostStorage(20),
				Comment: NewInMemoryCommentStorage(20),
			}

			result, err := s.Post.AddPost(test.input)
			if err != nil && !test.wantErr {
				t.Errorf("AddPost() error: %v, want error: %v", err, test.wantErr)
			}

			if !isPostsEqual(*result, *test.want) {
				t.Errorf("AddPost() result: %v, want: %v", result, test.want)
			}

		})

	}

}

func TestGetPostById(t *testing.T) {

	normalPost := model.Post{
		Author:          "normalAuthor",
		Header:          "Normal Header",
		Content:         "Normal content",
		CommentsAllowed: true,
	}

	s := InMemoryStorage{
		Post:    NewInMemoryPostStorage(20),
		Comment: NewInMemoryCommentStorage(20),
	}

	res, err := s.Post.AddPost(normalPost)
	if err != nil {
		t.Errorf("Can not add test post: %v", err)
	}

	tests := []struct {
		name    string
		input   int
		want    *model.Post
		wantErr bool
	}{
		{
			name:    "Normal post check",
			input:   res.ID,
			want:    &normalPost,
			wantErr: false,
		},
		{
			name:    "Id: -1 check",
			input:   -1,
			want:    &model.Post{},
			wantErr: true,
		},
		{
			name:    "Id: 4294967296 check",
			input:   4294967296,
			want:    &model.Post{},
			wantErr: true,
		},
		{
			name:    "Id: 0 check",
			input:   0,
			want:    &model.Post{},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			found, err := s.Post.GetPostById(test.input)
			if err != nil && !test.wantErr {
				t.Errorf("GetPostById() error: %v, want error: %v", err, test.wantErr)
			}

			if !isPostsEqual(found, *test.want) {
				t.Errorf("GetPostById() result: %v, want: %v", found, test.want)
			}

		})

	}

}
