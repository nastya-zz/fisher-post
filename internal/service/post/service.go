package post

import (
	"post/internal/client/db"
	userservice "post/internal/client/user_service"
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
