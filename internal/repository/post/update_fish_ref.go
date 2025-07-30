package post

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) UpdatePostFishReference(ctx context.Context, postId uuid.UUID, fishId int) error {
	const op = "repository.post.UpdateReferenceFish"

	builder := sq.Update(postFishTable).
		PlaceholderFormat(sq.Dollar).
		Set(fishIdColumn, fishId).
		Where(sq.Eq{postIdColumn: postId, fishIdColumn: fishId})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf(op+" error in update reference fish %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}
	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf(op+" error in update reference fish %w", err)
	}

	return nil
}
