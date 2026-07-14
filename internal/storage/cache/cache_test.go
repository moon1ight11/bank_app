package cache

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestCache(t *testing.T) (*CacheService, *miniredis.Miniredis) {
	t.Helper()

	mr, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	service := NewCacheService(client)
	return service, mr
}

func TestCache_SetAndGet(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "test:key"
	value := "hello world"

	err := cache.Set(ctx, key, value, 10*time.Minute)
	require.NoError(t, err)

	var result string
	err = cache.Get(ctx, key, &result)
	require.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestCache_Get_Miss(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	var result string
	err := cache.Get(ctx, "nonexistent", &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

func TestCache_SetAndGet_Struct(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	ctx := context.Background()
	original := testData{ID: 42, Name: "Тест"}

	err := cache.Set(ctx, "struct:key", original, 10*time.Minute)
	require.NoError(t, err)

	var result testData
	err = cache.Get(ctx, "struct:key", &result)
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestCache_Delete(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "delete:key"

	err := cache.Set(ctx, key, "some value", 10*time.Minute)
	require.NoError(t, err)

	err = cache.Delete(ctx, key)
	require.NoError(t, err)

	var result string
	err = cache.Get(ctx, key, &result)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cache miss")
}

func TestCache_Expiration(t *testing.T) {
	cache, mr := setupTestCache(t)
	defer mr.Close()

	ctx := context.Background()
	key := "expire:key"

	err := cache.Set(ctx, key, "will expire", 1*time.Second)
	require.NoError(t, err)

	mr.FastForward(2 * time.Second)

	var result string
	err = cache.Get(ctx, key, &result)
	require.Error(t, err)
}
