package comment

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
)

// Add implements repository.CommentRepository - создание комментария или ответа
func (r repo) Add(ctx context.Context, postID, userID uuid.UUID, content string, parentCommentID *uuid.UUID) (*model.Comment, error) {
	const op = "repository.comment.Add"

	// Если это ответ на комментарий, нужно получить ID пользователя, которому отвечают
	var replyToUserID *uuid.UUID
	if parentCommentID != nil {
		parentComment, err := r.getCommentByID(ctx, *parentCommentID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get parent comment: %w", op, err)
		}
		replyToUserID = &parentComment.UserID
	}

	builder := sq.Insert(commentsTable).
		Columns(postIdColumn, userIdColumn, contentColumn, parentCommentIdColumn, replyToUserIdColumn).
		Values(postID, userID, content, parentCommentID, replyToUserID).
		Suffix("RETURNING " + commentIdColumn + ", " + createdAtColumn + ", " + updatedAtColumn).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var result struct {
		CommentID uuid.UUID `db:"comment_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := r.db.DB().ScanOneContext(ctx, &result, q, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to create comment: %w", op, err)
	}

	// Возвращаем модель в формате API
	return &model.Comment{
		ID:              result.CommentID,
		PostID:          postID,
		UserID:          userID,
		ParentCommentID: parentCommentID,
		ReplyToUserID:   replyToUserID,
		Content:         content,
		CreatedAt:       result.CreatedAt,
		UpdatedAt:       result.UpdatedAt,
		IsReply:         parentCommentID != nil,
	}, nil
}
