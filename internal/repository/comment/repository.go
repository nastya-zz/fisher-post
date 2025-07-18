package comment

import (
	"context"

	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
	"post/internal/repository"
)

type repo struct {
	db db.Client
}

// GetCommentsCount implements repository.CommentRepository.
func (r repo) GetCommentsCount(ctx context.Context, postID uuid.UUID) (int, error) {
	panic("unimplemented")
}

func (r repo) Add(ctx context.Context, postID, userID uuid.UUID) (*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) Remove(ctx context.Context, postID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func New(db db.Client) repository.CommentRepository {
	return &repo{db: db}
}
