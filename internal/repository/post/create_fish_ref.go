package post

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

func (r repo) CreatePostFishReference(ctx context.Context, postId uuid.UUID, fishId int) error {
	const op = "repository.post.CreateReferenceFish"

	builder := sq.Insert(postFishTable).
		PlaceholderFormat(sq.Dollar).
		Columns(postIdColumn, fishIdColumn).
		Values(postId, fishId)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf(op+" error in create reference fish %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}
	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf(op+" error in create reference fish %w", err)
	}

	return nil
}
