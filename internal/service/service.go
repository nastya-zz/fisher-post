package service

import (
	"context"
	"github.com/google/uuid"
	"post/internal/model"
)

type PostService interface {
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
	AddLike(ctx context.Context, postID, userID uuid.UUID) (int, error)
	RemoveLike(ctx context.Context, postID, userID uuid.UUID) (int, error)
}

type CommentService interface {
	AddComment(ctx context.Context, postID, userID uuid.UUID) (*model.Comment, error)
	RemoveComment(ctx context.Context, postID, userID uuid.UUID) error
}
