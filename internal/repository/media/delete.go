package media

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.media.Delete"

	builder := sq.Delete(mediaTable).
		Where(sq.Eq{postIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf(op+" failed to build media query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf(op+" failed to delete media: %w", err)
	}

	return nil
}
