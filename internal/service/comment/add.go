package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	"post/pkg/logger"
)

// AddComment implements service.CommentService - создание комментария или ответа на комментарий
func (s serv) AddComment(ctx context.Context, postID, userID uuid.UUID, content string, parentCommentID *uuid.UUID) (*model.Comment, error) {
	const op = "comment.service.AddComment"

	// Валидация контента
	if content == "" {
		return nil, fmt.Errorf("%s: content cannot be empty", op)
	}

	if len(content) > 2000 {
		return nil, fmt.Errorf("%s: content too long (max 2000 characters)", op)
	}

	// Проверяем существование пользователя перед сохранением комментария
	_, err := s.userService.GetUser(ctx, "", userID)
	if err != nil {
		logger.Error("user does not exist", "error", err, "userID", userID)
		return nil, fmt.Errorf("%s: пользователь не существует: %w", op, err)
	}

	// Создаем комментарий через репозиторий
	comment, err := s.commentRepository.Add(ctx, postID, userID, content, parentCommentID)
	if err != nil {
		logger.Error("failed to add comment", "error", err, "postID", postID, "userID", userID)
		return nil, fmt.Errorf("%s: failed to add comment: %w", op, err)
	}

	// Обогащаем комментарий информацией о пользователе
	enrichedComments, err := s.enrichCommentsWithUsers(ctx, []*model.Comment{comment})
	if err != nil {
		logger.Error("failed to enrich comment with user", "error", err, "commentID", comment.ID)
		return nil, fmt.Errorf("%s: failed to enrich comment: %w", op, err)
	}

	if len(enrichedComments) == 0 {
		logger.Error("failed to enrich comment with user", "error", err, "commentID", comment.ID)
		return nil, fmt.Errorf("%s: не удалось обогатить комментарий данными пользователя", op)
	}

	logger.Info("comment added successfully", "commentID", comment.ID, "postID", postID, "userID", userID)

	return enrichedComments[0], nil
}
