package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/pkg/logger"
)

// RemoveComment implements service.CommentService - удаление комментария
func (s serv) RemoveComment(ctx context.Context, commentID, userID uuid.UUID) error {
	const op = "comment.service.RemoveComment"

	err := s.commentRepository.Remove(ctx, commentID, userID)
	if err != nil {
		logger.Error("failed to remove comment", "error", err, "commentID", commentID, "userID", userID)
		return fmt.Errorf("%s: failed to remove comment: %w", op, err)
	}

	logger.Info("comment removed successfully", "commentID", commentID, "userID", userID)
	return nil
}
