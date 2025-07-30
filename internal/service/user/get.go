package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/metadata"

	"post/internal/model"
)

func (s serv) GetUser(ctx context.Context, token string, id uuid.UUID) (*model.User, error) {
	md := metadata.Pairs("authorization", token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	profile, err := s.userClient.GetProfile(ctx, &desc.GetProfileRequest{
		Id: id.String(),
	})

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return profile, nil
}
