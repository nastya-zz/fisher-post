package repository

import (
	"context"
	"github.com/google/uuid"
	"post/internal/model"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	Update(ctx context.Context, post *model.Post) (*model.Post, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
	AddLike(ctx context.Context, postID, userID uuid.UUID) (int, error)
	RemoveLike(ctx context.Context, postID, userID uuid.UUID) (int, error)
}

type CommentRepository interface {
	Add(ctx context.Context, postID, userID uuid.UUID) (*model.Comment, error)
	Remove(ctx context.Context, postID, userID uuid.UUID) error
}
