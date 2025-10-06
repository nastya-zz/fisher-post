package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	"post/pkg/logger"
)

// GetReplies implements service.CommentService - получение ответов на конкретный комментарий
func (s serv) GetReplies(ctx context.Context, commentID uuid.UUID) ([]*model.Comment, error) {
	const op = "comment.service.GetReplies"

	replies, err := s.commentRepository.GetReplies(ctx, commentID)
	if err != nil {
		logger.Error("failed to get replies", "error", err, "commentID", commentID)
		return nil, fmt.Errorf("%s: failed to get replies: %w", op, err)
	}

	// Обогащаем ответы информацией о пользователях
	enrichedReplies, err := s.enrichCommentsWithUsers(ctx, replies)
	if err != nil {
		logger.Error("failed to enrich replies with users", "error", err, "commentID", commentID)
		return nil, fmt.Errorf("%s: failed to enrich replies: %w", op, err)
	}

	logger.Info("replies retrieved successfully", "commentID", commentID, "count", len(enrichedReplies))
	return enrichedReplies, nil
}
