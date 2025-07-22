package media

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
)

func (r *repo) Get(ctx context.Context, id uuid.UUID) (*model.Media, error) {
	const op = "repository.media.Get"

	builder := sq.Select(mediaIdColumn, mediaTypeColumn, mediaUrlColumn, mediaThumbnailUrlColumn).
		From(mediaTable).
		Where(sq.Eq{mediaIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf(op+" failed to build media query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var media model.Media
	err = r.db.DB().ScanOneContext(ctx, &media, q, args...)
	if err != nil {
		return nil, fmt.Errorf(op+" failed to get media: %w", err)
	}

	return &media, nil
}	