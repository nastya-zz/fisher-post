package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
	Ping(ctx context.Context) error
}

type CacheInvalidator interface {
	Invalidate(ctx context.Context, key string) error
}
