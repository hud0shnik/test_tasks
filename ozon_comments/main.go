package main

import (
	"comments/graph"
	"comments/internal/config"
	"comments/internal/storage"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const defaultPort = "8080"

func main() {

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	err := config.Init()
	if err != nil {
		log.Fatalf("InitConfig error: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var redisClient *redis.Client
	if os.Getenv(storage.STORAGE_FLAG) == "in-memory" {
		redisClient = storage.GetRedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0)
	} else {
		// Добавить коннект к постгрес
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Redis: redisClient,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
