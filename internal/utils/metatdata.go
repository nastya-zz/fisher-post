package utils

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"post/internal/model"
)

func GetUserIdFromMetadata(ctx context.Context) (uuid.UUID, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "метаданные отсутствуют")
	}

	userIDValues := md.Get("user-id")
	if len(userIDValues) == 0 {
		return uuid.Nil, status.Errorf(codes.Unauthenticated, "user-id не найден в заголовках")
	}

	userID, err := model.GetUuid(userIDValues[0])
	if err != nil {
		return uuid.Nil, status.Errorf(codes.InvalidArgument, "Id пользователя не валидный")
	}

	return userID, nil
}
