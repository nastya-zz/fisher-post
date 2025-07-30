package casheduserservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"post/internal/model"
)

func (s serv) GetUser(ctx context.Context, token string, id uuid.UUID) (*model.User, error) {
	const op = "service.cached_user.GetUser"

	// Формируем ключ кеша
	cacheKey := fmt.Sprintf("user:%s", id.String())

	// Пытаемся получить из кеша
	var cachedUser model.User
	if err := s.cacheService.Get(ctx, cacheKey, &cachedUser); err == nil {
		return &cachedUser, nil
	}

	// При промахе - идем в основной сервис
	user, err := s.userService.GetUser(ctx, token, id)
	if err != nil {
		// Если пользователь не найден - удаляем из кеша
		if isNotFoundError(err) {
			s.cacheService.Delete(ctx, cacheKey)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Кешируем успешный результат
	if err := s.cacheService.Set(ctx, cacheKey, user, s.ttl); err != nil {
		// Логируем ошибку кеширования, но не падаем
		// logger.Warn("failed to cache user", "error", err)
	}

	return user, nil
}

func isNotFoundError(err error) bool {
	// Проверяем на gRPC NotFound или другие коды "не найден"
	return strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "NotFound")
}
