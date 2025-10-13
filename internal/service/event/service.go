package event

import (
	"post/internal/producer/rabbitmq"
	"post/internal/repository"
	"post/internal/service"
)

type EventService struct {
	cachedUserService service.CachedUserService
	postService service.PostService
	producer rabbitmq.Producer
	eventRepository repository.EventRepository
}

func New(cachedUserService service.CachedUserService, postService service.PostService, producer rabbitmq.Producer, eventRepository repository.EventRepository) service.EventsService {
	return &EventService{cachedUserService: cachedUserService, postService: postService, producer: producer, eventRepository: eventRepository}
}
