package post

import (
	"context"

	"github.com/google/uuid"

	"post/internal/client/db"
	userservice "post/internal/client/user_service"
	"post/internal/model"
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	userService userservice.ServiceClient
	repository  repository.PostRepository
	likeRepository repository.LikeRepository
	txManager   db.TxManager
}

func (s serv) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
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

func New(repository repository.PostRepository, manager db.TxManager, userService userservice.ServiceClient, likeRepository repository.LikeRepository) service.PostService {
	return &serv{
		repository:  repository,
		txManager:   manager,
		userService: userService,
		likeRepository: likeRepository,
	}
}
