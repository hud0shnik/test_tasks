package graph

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Redis   *redis.Client
	Posgres *sql.DB
}
