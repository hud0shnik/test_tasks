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

	var postgresDB *sqlx.DB
	var inMemoryStorage *storage.InMemoryStorage

	// Проверка на место хранения данных
	if os.Getenv(storage.STORAGE_FLAG) == "db" {
		log.Info("Connecting to Postgres...")
		postgresDB, err = storage.GetPostgresDB()
		if err != nil {
			log.Fatalf("GetPostgresDB error: %v", err)
		}
	} else {
		log.Info("Creating in-memory storage...")
		inMemoryStorage = &storage.InMemoryStorage{
			Post:    storage.NewInMemoryPostStorage(10),
			Comment: storage.NewInMemoryCommentStorage(10),
		}

	}

	// Создание сервера
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			Posgres:         postgresDB,
			InMemoryStorage: inMemoryStorage,
		}}))

	// Ручки сервиса
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Запуск сервера
	log.Infof("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
