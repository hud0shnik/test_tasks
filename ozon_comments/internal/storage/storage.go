package storage

import (
	"comments/graph/model"
	"errors"
	"os"
)

const STORAGE_FLAG = "STORAGE"

func SetPost(post model.Post) error {
	if os.Getenv(STORAGE_FLAG) == "in-memory" {
		return SetPostInMemory(post)
	} else {
		return errors.New("SetPostInDB not implemented")
	}

}

func SetPostInMemory(model.Post) error {
	return nil
}
