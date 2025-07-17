package userservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Исправлено форматирование импорта: добавлена пустая строка после блока импортов для соответствия линтеру

const servicePort = 50052

type Client struct {
	ctx context.Context
	md  metadata.MD
}
type userService struct {
	cl desc.UserV1Client
}

type ServiceClient interface {
	GetUser(ctx context.Context, token string, id uuid.UUID) (*desc.GetProfileResponse, error)
}

func New(ctx context.Context) (ServiceClient, error) {

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", servicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial GRPC client: %w", err)
	}
	defer func() { _ = conn.Close() }()

	cl := desc.NewUserV1Client(conn)

	return &userService{
		cl: cl,
	}, nil
}

func (userService userService) GetUser(ctx context.Context, token string, id uuid.UUID) (*desc.GetProfileResponse, error) {
	md := metadata.Pairs("authorization", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	profile, err := userService.cl.GetProfile(ctx, &desc.GetProfileRequest{
		Id: id.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return profile, nil
}
