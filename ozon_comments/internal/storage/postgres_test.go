package storage

import (
	"comments/graph/model"
	"testing"

	"github.com/joho/godotenv"
)

func TestSavePostToPostgres(t *testing.T) {

	godotenv.Load("../../.env")

	db, err := GetPostgresDB()
	if err != nil {
		t.Errorf("can not connect to Postgres: %v", err)
	}

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

			result, err := SavePostToPostgres(db, test.input)
			if err != nil && !test.wantErr {
				t.Errorf("AddPost() error: %v, want error: %v", err, test.wantErr)
			}

			if !isPostsEqual(*result, *test.want) {
				t.Errorf("AddPost() result: %v, want: %v", result, test.want)
			}

		})

	}

}
