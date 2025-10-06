package comment

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
)

// GetReplies implements repository.CommentRepository - получение ответов на конкретный комментарий
func (r repo) GetReplies(ctx context.Context, commentID uuid.UUID) ([]*model.Comment, error) {
	const op = "repository.comment.GetReplies"

	// Получаем ответы на комментарий
	repliesQuery := sq.Select("*").
		From(commentsTable).
		Where(sq.Eq{parentCommentIdColumn: commentID}).
		OrderBy(createdAtColumn + " ASC").
		PlaceholderFormat(sq.Dollar)

	repliesSql, repliesArgs, err := repliesQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build replies query: %w", op, err)
	}

	repliesQ := db.Query{
		Name:     op,
		QueryRaw: repliesSql,
	}

	var replies []Comment
	if err := r.db.DB().ScanAllContext(ctx, &replies, repliesQ, repliesArgs...); err != nil {
		return nil, fmt.Errorf("%s: failed to get replies: %w", op, err)
	}

	// Конвертируем в API модель
	var result []*model.Comment
	for _, reply := range replies {
		result = append(result, r.convertToModelComment(&reply))
	}

	return result, nil
}
