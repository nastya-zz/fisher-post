package event

import (
	"post/internal/service"
)

type Processor struct {
	cachedUserService service.CachedUserService
	postService service.PostService
}

func New(cachedUserService service.CachedUserService, postService service.PostService) service.EventsService {
	return &Processor{cachedUserService: cachedUserService, postService: postService}
}
