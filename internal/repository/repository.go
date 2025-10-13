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
	DeleteByUserId(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type LikeRepository interface {
	Add(ctx context.Context, postID, userID uuid.UUID) error
	RemoveLike(ctx context.Context, postID, userID uuid.UUID) error
	GetLikesCount(ctx context.Context, postID uuid.UUID) (int, error)
	GetLikes(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) //список ид пользователей, которые лайкнули пост
}

type CommentRepository interface {
	// Создание комментария или ответа на комментарий
	Add(ctx context.Context, postID, userID uuid.UUID, content string, parentCommentID *uuid.UUID) (*model.Comment, error)

	// Получение всех комментариев к посту с вложенностью (один уровень)
	GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error)

	// Получение ответов на конкретный комментарий
	GetReplies(ctx context.Context, commentID uuid.UUID) ([]*model.Comment, error)

	// Удаление комментария (только своим автором)
	Remove(ctx context.Context, commentID, userID uuid.UUID) error

	// Подсчет общего количества комментариев к посту (включая ответы)
	GetCommentsCount(ctx context.Context, postID uuid.UUID) (int, error)
}

type MediaRepository interface {
	Create(ctx context.Context, media *model.CreateMedia) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*model.Media, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByPostID(ctx context.Context, postID uuid.UUID) ([]model.Media, error)
}

type EventRepository interface {
	GetNewEvent(ctx context.Context, batchSize int) ([]*model.Event, error)
	SaveEvent(ctx context.Context, event *model.Event) error
	SetDone(ctx context.Context, id int) error
}
