package graph

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Redis   *redis.Client
	Posgres *sqlx.DB
}
