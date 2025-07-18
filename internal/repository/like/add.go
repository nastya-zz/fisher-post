package like

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) Add(ctx context.Context, postID, userID uuid.UUID) error {
	const op = "repository.like.Add"

	builder := sq.Insert(likesTable).
		Columns(postIdColumn, userIdColumn).
		Values(postID, userID).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf(op+" failed to build query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf(op+" failed to add like: %w", err)
	}

	return nil
}
