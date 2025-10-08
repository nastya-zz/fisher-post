package userservice

import (
	"context"
	"fmt"
	"os"

	desc "github.com/nastya-zz/fisher-protocols/gen/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"post/pkg/logger"
)

const servicePort = 50052

type userService struct {
	Cl   desc.UserV1Client
	conn *grpc.ClientConn
}

type ServiceClient interface {
	Close() error
	GetProfile(ctx context.Context, req *desc.GetProfileRequest) (*desc.GetProfileResponse, error)
}

func New(ctx context.Context) (ServiceClient, error) {
	host := mustUserServiceHost("localhost")

	logger.Info("dialing user service", "host", host)
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%d", host, servicePort),
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

func mustUserServiceHost(defaultVal string) string {
	host := os.Getenv("USER_SERVICE_HOST")
	if host == "" {
		return defaultVal
	}

	return host
}
