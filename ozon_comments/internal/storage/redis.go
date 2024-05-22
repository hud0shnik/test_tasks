package storage

import (
	"comments/graph/model"

	"github.com/redis/go-redis/v9"
)

func GetRedisClient(addr, password string, db int) *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return redisClient
}

func SavePostToRedis(db *redis.Client, post model.Post) error {

	// Добавить сохранение в Redis

	return nil

}
