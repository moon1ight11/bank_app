package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// метод для установки кэша
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// маршалим данные, которые будем кэшировать
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// кэшируем
	return c.client.Set(ctx, key, data, expiration).Err()
}

// метод для получения кэша
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	// получаем данные и кэша
	data, err := c.client.Get(ctx, key).Bytes()

	// проверяем ошибки
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss: %w", err)
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	return json.Unmarshal(data, dest)
}

// метод для удаления данных из кэша
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
