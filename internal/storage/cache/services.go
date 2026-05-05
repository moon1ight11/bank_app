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
		return fmt.Errorf("Error in set cache: %w", err)
	}

	// кэшируем
	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("Error in set cache: %w", err)
	}

	return nil
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
		return fmt.Errorf("Error in get cache: %w", err)
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return fmt.Errorf("Error in get cache: %w", err)
	}

	return nil
}

// метод для удаления данных из кэша
func (c *CacheService) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("Error in delete cache: %w", err)
	}

	return nil
}
