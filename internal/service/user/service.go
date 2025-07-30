package user

import (
	userservice "post/internal/client/user_service"
	"post/internal/service"
)

type serv struct {
	userClient userservice.ServiceClient
}

func New(userClient userservice.ServiceClient) service.UserService {
	return &serv{
		userClient: userClient,
	}
}
