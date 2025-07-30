package repository

import (
	"context"

	"github.com/google/uuid"

	"post/internal/model"
	repoModel "post/internal/repository/post/model"
)

type PostRepository interface {
	CreatePost(ctx context.Context, post *model.CreatePost) (*repoModel.CreatedPost, error)
	CreatePostFishReference(ctx context.Context, postId uuid.UUID, fishId int) error
	CreatePostTackleReference(ctx context.Context, postId uuid.UUID, tackleId int) error
	UpdatePostFishReference(ctx context.Context, postId uuid.UUID, fishId int) error
	UpdatePostTackleReference(ctx context.Context, postId uuid.UUID, tackleId int) error

	Update(ctx context.Context, post *model.UpdatePost) error
	Get(ctx context.Context, id uuid.UUID) (*repoModel.Post, error)
	GetPosts(ctx context.Context, id uuid.UUID) ([]*repoModel.Post, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type LikeRepository interface {
	Add(ctx context.Context, postID, userID uuid.UUID) error
	RemoveLike(ctx context.Context, postID, userID uuid.UUID) error
	GetLikesCount(ctx context.Context, postID uuid.UUID) (int, error)
	GetLikes(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) //список ид пользователей, которые лайкнули пост
}

type CommentRepository interface {
	Add(ctx context.Context, postID, userID uuid.UUID) (*model.Comment, error)
	Remove(ctx context.Context, postID, userID uuid.UUID) error
	GetCommentsCount(ctx context.Context, postID uuid.UUID) (int, error)
}

type MediaRepository interface {
	Create(ctx context.Context, media *model.CreateMedia) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Media, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByPostID(ctx context.Context, postID uuid.UUID) ([]model.Media, error)
}
