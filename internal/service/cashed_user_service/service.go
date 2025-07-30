package casheduserservice

import (
	"context"
	"time"

	"post/internal/client/cache"
	"post/internal/service"
)

type serv struct {
	userService  service.UserService
	cacheService cache.Cache
	invalidator  cache.CacheInvalidator
	ttl          time.Duration
}

func New(userService service.UserService, cacheService cache.Cache, invalidator cache.CacheInvalidator) service.CachedUserService {
	return &serv{
		userService:  userService,
		cacheService: cacheService,
		invalidator:  invalidator,
		ttl:          15 * time.Minute,
	}
}

func (s serv) Invalidate(ctx context.Context, key string) error {
	return s.invalidator.Invalidate(ctx, key)
}