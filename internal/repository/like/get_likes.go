package like

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) GetLikes(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) {
	const op = "repository.like.GetLikes"

	builder := sq.Select(userIdColumn).
		From(likesTable).
		Where(sq.Eq{postIdColumn: postID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf(op+" failed to build query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var ids []uuid.UUID
	err = r.db.DB().ScanAllContext(ctx, &ids, q, args...)
	if err != nil {
		return nil, fmt.Errorf(op+" failed to get likes: %w", err)
	}

	return ids, nil
}
