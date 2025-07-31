package userservice

import (
	"context"
	"fmt"

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
	Cl   desc.UserV1Client
	conn *grpc.ClientConn
}

type ServiceClient interface {
	Close() error
	GetProfile(ctx context.Context, req *desc.GetProfileRequest) (*desc.GetProfileResponse, error)
}

func New(ctx context.Context) (ServiceClient, error) {

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", servicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial GRPC client: %w", err)
	}

	cl := desc.NewUserV1Client(conn)

	return &userService{
		Cl:   cl,
		conn: conn,
	}, nil
}
func (s userService) GetProfile(ctx context.Context, req *desc.GetProfileRequest) (*desc.GetProfileResponse, error) {
	return s.Cl.GetProfile(ctx, req)
}
func (userService userService) Close() error {
	if userService.Cl != nil {
		userService.conn.Close()
	}

	return nil
}
