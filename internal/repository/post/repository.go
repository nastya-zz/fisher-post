package post

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

func (r repo) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) Update(ctx context.Context, post *model.Post) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) Get(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r repo) AddLike(ctx context.Context, postID, userID uuid.UUID) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) RemoveLike(ctx context.Context, postID, userID uuid.UUID) (int, error) {
	//TODO implement me
	panic("implement me")
}

func New(db db.Client) repository.PostRepository {
	return &repo{db: db}
}
