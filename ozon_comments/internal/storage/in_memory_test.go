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
