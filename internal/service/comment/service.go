package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	"post/internal/repository"
	"post/internal/service"
	"post/pkg/logger"
)

type serv struct {
	commentRepository repository.CommentRepository
	userService       service.CachedUserService
}

// enrichCommentsWithUsers - обогащает комментарии информацией о пользователях
// Возвращает только комментарии, для которых удалось получить данные пользователя
func (s serv) enrichCommentsWithUsers(ctx context.Context, comments []*model.Comment) ([]*model.Comment, error) {
	var enrichedComments []*model.Comment

	for _, comment := range comments {
		// Пропускаем комментарии без указанного пользователя
		if comment.UserID == uuid.Nil {
			logger.Debug("skipping comment with nil userID", "commentID", comment.ID)
			continue
		}

		// Получаем данные пользователя
		user, err := s.userService.GetUser(ctx, "", comment.UserID)
		if err != nil {
			logger.Error("failed to get user for comment", "error", err, "userID", comment.UserID)
			// Пропускаем комментарий, если не удалось получить данные пользователя
			continue
		}

		comment.User = user

		// Рекурсивно обогащаем ответы
		if len(comment.Replies) > 0 {
			enrichedReplies, err := s.enrichCommentsWithUsers(ctx, comment.Replies)
			if err != nil {
				return nil, fmt.Errorf("failed to enrich replies for comment %s: %w", comment.ID, err)
			}
			comment.Replies = enrichedReplies
		}

		// Добавляем успешно обогащенный комментарий в результат
		enrichedComments = append(enrichedComments, comment)
	}

	return enrichedComments, nil
}

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

func New(commentRepository repository.CommentRepository, userService service.CachedUserService) service.CommentService {
	return &serv{
		commentRepository: commentRepository,
		userService:       userService,
	}
}
