package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	"post/pkg/logger"
)

// GetCommentsByPostID implements service.CommentService - получение всех комментариев к посту с вложенностью
func (s serv) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	const op = "comment.service.GetCommentsByPostID"

	// Валидация параметров пагинации
	if limit <= 0 || limit > 100 {
		limit = 20 // дефолтное значение
	}
	if offset < 0 {
		offset = 0
	}

	comments, err := s.commentRepository.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		logger.Error("failed to get comments", "error", err, "postID", postID)
		return nil, fmt.Errorf("%s: failed to get comments: %w", op, err)
	}

	// Обогащаем комментарии информацией о пользователях
	enrichedComments, err := s.enrichCommentsWithUsers(ctx, comments)
	if err != nil {
		logger.Error("failed to enrich comments with users", "error", err, "postID", postID)
		return nil, fmt.Errorf("%s: failed to enrich comments: %w", op, err)
	}

	logger.Info("comments retrieved successfully", "postID", postID, "count", len(enrichedComments))
	return enrichedComments, nil
}

// GetCommentsCount implements service.CommentService - подсчет общего количества комментариев к посту
func (s serv) GetCommentsCount(ctx context.Context, postID uuid.UUID) (int, error) {
	const op = "comment.service.GetCommentsCount"

	count, err := s.commentRepository.GetCommentsCount(ctx, postID)
	if err != nil {
		logger.Error("failed to get comments count", "error", err, "postID", postID)
		return 0, fmt.Errorf("%s: failed to get comments count: %w", op, err)
	}

	logger.Info("comments count retrieved successfully", "postID", postID, "count", count)
	return count, nil
}
