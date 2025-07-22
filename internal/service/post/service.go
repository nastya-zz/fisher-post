package post

import (
	"context"

	"post/internal/client/db"
	userservice "post/internal/client/user_service"
	"post/internal/model"
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	userService     userservice.ServiceClient
	repository      repository.PostRepository
	likeRepository  repository.LikeRepository
	txManager       db.TxManager
	mediaRepository repository.MediaRepository
	minioService    service.MinioService
}

func (s serv) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func New(repository repository.PostRepository,
	manager db.TxManager,
	userService userservice.ServiceClient,
	likeRepository repository.LikeRepository,
	mediaRepository repository.MediaRepository,
	minioService service.MinioService,
) service.PostService {
	return &serv{
		repository:      repository,
		txManager:       manager,
		userService:     userService,
		likeRepository:  likeRepository,
		mediaRepository: mediaRepository,
		minioService:    minioService,
	}
}
