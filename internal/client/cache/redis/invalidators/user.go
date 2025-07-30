package invalidators

import (
	"context"

	"post/internal/client/cache"
	"post/pkg/logger"
)

type cacheInvalidator struct {
	cache cache.Cache
}

func NewCacheInvalidator(cache cache.Cache) cache.CacheInvalidator {
	return &cacheInvalidator{cache: cache}
}

func (i *cacheInvalidator) Invalidate(ctx context.Context, key string) error {
	logger.Info("Invalidating user cache", "User.ID ", key)
	return i.cache.Delete(ctx, key)
}
