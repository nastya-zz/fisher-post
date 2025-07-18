package like

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r *repo) GetLikesCount(ctx context.Context, postID uuid.UUID) (int, error) {
	const op = "repository.like.GetLikesCount"

	builder := sq.Select("count(*)").
		From(likesTable).
		Where(sq.Eq{postIdColumn: postID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf(op+" failed to build query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var count int
	if err := r.db.DB().QueryRowContext(ctx, q, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf(op+" failed to get likes count: %w", err)
	}

	return count, nil
}
