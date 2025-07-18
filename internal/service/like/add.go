package like

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s serv) AddLike(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (int, error) {
	const op = "service.like.AddLike"

	if err := s.repository.Add(ctx, postID, userID); err != nil {
		return 0, fmt.Errorf(op+" failed to add like: %w", err)
	}

	likesCount, err := s.repository.GetLikesCount(ctx, postID)
	if err != nil {
		return 0, fmt.Errorf(op+" failed to get likes count: %w", err)
	}

	return likesCount, nil
}
