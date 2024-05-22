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
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

const defaultPort = "8080"

func main() {

	// Настройка логгера
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.DateTime,
	})

	// Загрузка конфига (переменных окружения)
	err := config.Init()
	if err != nil {
		log.Fatalf("InitConfig error: %v", err)
	}

	// Получение порта для плейграунда
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	var redisClient *redis.Client
	var postgresDB *sqlx.DB

	// Проверка на место хранения данных
	if os.Getenv(storage.STORAGE_FLAG) == "in-memory" {
		log.Info("Connecting to Redis...")
		redisClient = storage.GetRedisClient(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0)
	} else {
		log.Info("Connecting to Postgres...")
		postgresDB, err = storage.GetPostgresDB()
		if err != nil {
			log.Fatalf("GetPostgresDB error: %v", err)
		}
	}

	// Создание сервера
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			Redis:   redisClient,
			Posgres: postgresDB,
		}}))

	// Ручки сервиса
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Запуск сервера
	log.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
