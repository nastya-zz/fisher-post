package post

import (
	"context"
	"strings"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/converter"
	"post/internal/model"
	"post/pkg/logger"
)

// AddComment добавляет комментарий к посту
func (i *Implementation) AddComment(ctx context.Context, req *desc.AddCommentRequest) (*desc.Comment, error) {
	const op = "api.post.AddComment"

	// Валидация запроса
	errors := i.validateAddCommentRequest(req)
	if len(errors) > 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s", strings.Join(errors, ", "))
	}

	// Парсим ID поста
	postID, err := model.GetUuid(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid post id", op)
	}

	// Парсим ID пользователя
	userID, err := model.GetUuid(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid user id", op)
	}

	// Создаем комментарий через сервис (пока без поддержки вложенных комментариев)
	comment, err := i.commentService.AddComment(ctx, postID, userID, req.Content, nil)
	if err != nil {
		logger.Error("failed to add comment", "error", err, "postID", postID, "userID", userID)
		return nil, status.Errorf(codes.Internal, "%s: failed to add comment", op)
	}

	logger.Info("comment added successfully", "commentID", comment.ID, "postID", postID, "userID", userID)
	return converter.FromModelCommentToDescComment(comment), nil
}

// GetComments получает все комментарии к посту
func (i *Implementation) GetComments(ctx context.Context, req *desc.PostId) (*desc.GetCommentsResponse, error) {
	const op = "api.post.GetComments"

	// Валидация запроса
	if req.GetPostId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "%s: post id is required", op)
	}

	// Парсим ID поста
	postID, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid post id: %v", op, err)
	}

	// Получаем комментарии через сервис (используем дефолтную пагинацию)
	comments, err := i.commentService.GetCommentsByPostID(ctx, postID, 20, 0)
	if err != nil {
		logger.Error("failed to get comments", "error", err, "postID", postID)
		return nil, status.Errorf(codes.Internal, "%s: failed to get comments", op)
	}

	// Конвертируем в protobuf формат
	pbComments := converter.FromModelCommentsToDescComments(comments)

	logger.Info("comments retrieved successfully", "postID", postID, "count", len(pbComments))
	return &desc.GetCommentsResponse{
		Comments: pbComments,
	}, nil
}

// RemoveComment удаляет комментарий
func (i *Implementation) RemoveComment(ctx context.Context, req *desc.RemoveCommentRequest) (*desc.RemoveCommentResponse, error) {
	const op = "api.post.RemoveComment"

	// Валидация запроса
	errors := i.validateRemoveCommentRequest(req)
	if len(errors) > 0 {
		return nil, status.Errorf(codes.InvalidArgument, "%s", strings.Join(errors, ", "))
	}

	// Парсим ID комментария
	commentID, err := model.GetUuid(req.CommentId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid comment id", op)
	}

	// Парсим ID пользователя
	userID, err := model.GetUuid(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s: invalid user id", op)
	}

	// Удаляем комментарий через сервис
	err = i.commentService.RemoveComment(ctx, commentID, userID)
	if err != nil {
		logger.Error("failed to remove comment", "error", err, "commentID", commentID, "userID", userID)
		return nil, status.Errorf(codes.Internal, "%s: failed to remove comment", op)
	}

	logger.Info("comment removed successfully", "commentID", commentID, "userID", userID)
	return &desc.RemoveCommentResponse{
		Success: true,
	}, nil
}

// validateAddCommentRequest валидирует запрос на добавление комментария
func (i *Implementation) validateAddCommentRequest(req *desc.AddCommentRequest) []string {
	var errors []string

	if req.GetPostId() == "" {
		errors = append(errors, "post id is required")
	}

	if req.GetUserId() == "" {
		errors = append(errors, "user id is required")
	}

	if req.GetContent() == "" {
		errors = append(errors, "content cannot be empty")
	}

	if len(req.GetContent()) > 2000 {
		errors = append(errors, "content too long (max 2000 characters)")
	}

	return errors
}

// validateRemoveCommentRequest валидирует запрос на удаление комментария
func (i *Implementation) validateRemoveCommentRequest(req *desc.RemoveCommentRequest) []string {
	var errors []string

	if req.GetCommentId() == "" {
		errors = append(errors, "comment id is required")
	}

	if req.GetUserId() == "" {
		errors = append(errors, "user id is required")
	}

	return errors
}
