package comment

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx"

	"post/internal/client/db"
)

// getCommentByID - вспомогательный метод для получения комментария по ID
func (r repo) getCommentByID(ctx context.Context, commentID uuid.UUID) (*Comment, error) {
	const op = "repository.comment.getCommentByID"

	builder := sq.Select("*").
		From(commentsTable).
		Where(sq.Eq{commentIdColumn: commentID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var comment Comment
	if err := r.db.DB().ScanOneContext(ctx, &comment, q, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: comment not found", op)
		}
		return nil, fmt.Errorf("%s: failed to get comment: %w", op, err)
	}

	return &comment, nil
}
