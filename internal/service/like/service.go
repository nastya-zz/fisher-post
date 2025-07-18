package like

import (
	"context"

	"github.com/google/uuid"

	"post/internal/model"
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	repository repository.LikeRepository
}

// GetLikes implements service.LikeService.
func (s *serv) GetLikes(ctx context.Context, postID uuid.UUID) ([]model.User, error) {
	panic("unimplemented")
}

// RemoveLike implements service.LikeService.
func (s *serv) RemoveLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (int, error) {
	panic("unimplemented")
}

func New(repository repository.LikeRepository) service.LikeService {
	return &serv{repository: repository}
}
