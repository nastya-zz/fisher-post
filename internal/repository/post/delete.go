package post

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.post.Delete"

	builder := sq.Delete(postsTable).
		Where(sq.Eq{postIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf(op+" failed to build post query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf(op+" failed to delete post: %w", err)
	}

	return nil
}
