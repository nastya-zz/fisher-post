package comment

import (
	"post/internal/repository"
	"post/internal/service"
)

type serv struct {
	commentRepository repository.CommentRepository
	userService       service.CachedUserService
}

func New(commentRepository repository.CommentRepository, userService service.CachedUserService) service.CommentService {
	return &serv{
		commentRepository: commentRepository,
		userService:       userService,
	}
}
