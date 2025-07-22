package like

import (
	"post/internal/client/user_service"
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	repository  repository.LikeRepository
	userService userservice.ServiceClient
}

// GetLikes implements service.LikeService.

func New(repository repository.LikeRepository, userService userservice.ServiceClient) service.LikeService {
	return &serv{repository: repository, userService: userService}
}
