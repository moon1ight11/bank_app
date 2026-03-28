package storage

import (
	"bank_app/internal/config"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// соединение с редис
func NewRedisClient(cfg *config.Config) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

// метод для закрытия соединения с редис
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
