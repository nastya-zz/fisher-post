package service

import (
	"context"

	"github.com/google/uuid"

	"post/internal/client/cache"
	"post/internal/model"
)

type PostService interface {
	CreatePost(ctx context.Context, post *model.CreatePost) (*model.Post, error)
	UpdatePost(ctx context.Context, post *model.UpdatePost) (*model.Post, error)
	GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error)
	GetPosts(ctx context.Context, id uuid.UUID) ([]*model.Post, error)
	DeletePost(ctx context.Context, id uuid.UUID) error
	UploadMedia(ctx context.Context, media *model.DescCreateMedia) (uuid.UUID, error)
	DeleteByUserId(ctx context.Context, userId uuid.UUID) error
}

type CommentService interface {
	// AddComment Создание комментария или ответа на комментарий
	AddComment(ctx context.Context, postID, userID uuid.UUID, content string, parentCommentID *uuid.UUID) (*model.Comment, error)

	// GetCommentsByPostID Получение всех комментариев к посту с вложенностью (один уровень)
	GetCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error)

	// GetReplies Получение ответов на конкретный комментарий
	GetReplies(ctx context.Context, commentID uuid.UUID) ([]*model.Comment, error)

	// RemoveComment Удаление комментария (только своим автором)
	RemoveComment(ctx context.Context, commentID, userID uuid.UUID) error

	// GetCommentsCount Подсчет общего количества комментариев к посту (включая ответы)
	GetCommentsCount(ctx context.Context, postID uuid.UUID) (int, error)
}

type LikeService interface {
	GetLikes(ctx context.Context, postID uuid.UUID) ([]*model.User, error)
	AddLike(ctx context.Context, postID, userID uuid.UUID) (int, error)
	RemoveLike(ctx context.Context, postID, userID uuid.UUID) (int, error)
}

type MinioService interface {
	UploadFile(ctx context.Context, file []byte, filename string) (string, error)
	RemoveFile(ctx context.Context, link string) error
}

type UserService interface {
	GetUser(ctx context.Context, token string, id uuid.UUID) (*model.User, error)
}

type CachedUserService interface {
	UserService
	cache.CacheInvalidator
}

type EventsService interface {
	Process(ctx context.Context, event model.Event) error
}
