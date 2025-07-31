package like

import (
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	repository  repository.LikeRepository
	userService service.UserService
}

// GetLikes implements service.LikeService.

func New(repository repository.LikeRepository, userService service.UserService) service.LikeService {
	return &serv{repository: repository, userService: userService}
}
