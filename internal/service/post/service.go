package post

import (
	"context"
	"github.com/google/uuid"
	"post/internal/client/db"
	"post/internal/model"
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	repository repository.PostRepository
	txManager  db.TxManager
}

func (s serv) CreatePost(ctx context.Context, post *model.CreatePost) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s serv) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s serv) GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s serv) DeletePost(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s serv) AddLike(ctx context.Context, postID, userID uuid.UUID) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (s serv) RemoveLike(ctx context.Context, postID, userID uuid.UUID) (int, error) {
	//TODO implement me
	panic("implement me")
}

func New(repository repository.PostRepository, manager db.TxManager) service.PostService {
	return &serv{
		repository: repository,
		txManager:  manager,
	}
}
