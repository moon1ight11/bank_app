package storage

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
)

type DataBase struct {
	DB            *sql.DB
	MigrationsDir string
}

// структура клиента редис
type RedisClient struct {
	Client *redis.Client
}
