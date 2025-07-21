package like

import (
	"context"

	"github.com/google/uuid"
)

func (s serv) RemoveLike(ctx context.Context, postID, userID uuid.UUID) (int, error) {

	if err := s.repository.RemoveLike(ctx, postID, userID); err != nil {
		return 0, err
	}

	count, err := s.repository.GetLikesCount(ctx, postID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
