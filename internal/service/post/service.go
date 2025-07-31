package post

import (
	"post/internal/client/db"
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	userService     service.CachedUserService
	repository      repository.PostRepository
	likeRepository  repository.LikeRepository
	txManager       db.TxManager
	mediaRepository repository.MediaRepository
	minioService    service.MinioService
}

func New(repository repository.PostRepository,
	manager db.TxManager,
	userService service.CachedUserService,
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
