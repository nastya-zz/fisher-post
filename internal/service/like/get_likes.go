package like

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
)

func (s serv) GetLikes(ctx context.Context, postID uuid.UUID) ([]*model.User, error) {
	const op = "service.like.GetLikes"

	userIds, err := s.repository.GetLikes(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf(op+" failed to get likes: %w", err)
	}

	var users []*model.User

	for _, userId := range userIds {
		user, err := s.userService.GetUser(ctx, "", userId)
		if err != nil {
			return nil, fmt.Errorf(op+" failed to get user: %w", err)
		}
		users = append(users, &model.User{
			ID:        user.ID,
			Username:  user.Username,
			AvatarUrl: user.AvatarUrl,
		})
	}

	return users, nil
}
