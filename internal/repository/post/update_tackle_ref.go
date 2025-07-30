package post

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) UpdatePostTackleReference(ctx context.Context, postId uuid.UUID, tackleId int) error {
	const op = "repository.post.UpdateReferenceTackle"

	builder := sq.Update(postTackleTable).
		PlaceholderFormat(sq.Dollar).
		Set(tackleIdColumn, tackleId).
		Where(sq.Eq{postIdColumn: postId, tackleIdColumn: tackleId})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf(op+" error in update reference tackle %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}
	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf(op+" error in update reference tackle %w", err)
	}

	return nil
}
