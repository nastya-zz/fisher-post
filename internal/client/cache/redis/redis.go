package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"

	"post/internal/client/cache"
	"post/pkg/logger"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) cache.Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &redisCache{client: rdb}
}

func (r *redisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *redisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}

func (r *redisCache) Ping(ctx context.Context) error {
	pong, err := r.client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	logger.Info("redis started successfully", "ping", pong)
	return nil
}
