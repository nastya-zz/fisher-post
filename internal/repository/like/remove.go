package like

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) Remove(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	const op = "repository.like.Remove"

	builder := sq.Delete(likesTable).
		Where(sq.Eq{postIdColumn: postID, userIdColumn: userID}).
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
		return fmt.Errorf(op+" failed to remove like: %w", err)
	}

	return nil
}
