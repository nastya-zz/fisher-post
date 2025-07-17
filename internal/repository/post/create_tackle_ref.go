package post

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) CreatePostTackleReference(ctx context.Context, postId uuid.UUID, tackleId int) error {
	const op = "repository.post.CreatePostTackleReference"

	builder := sq.Insert(postTackleTable).
		PlaceholderFormat(sq.Dollar).
		Columns(postIdColumn, tackleIdColumn).
		Values(postId, tackleId)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf(op+" error in create reference tackle %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}
	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf(op+" error in create reference tackle %w", err)
	}

	return nil
}
