package media

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
	"post/pkg/logger"
)

func (r *repo) Create(ctx context.Context, media *model.CreateMedia) (uuid.UUID, error) {
	const op = "repository.media.Create"

	builder := sq.Insert(mediaTable).
		PlaceholderFormat(sq.Dollar).
		Columns(
			postIdColumn,
			mediaTypeColumn,
			mediaUrlColumn,
			mediaThumbnailUrlColumn,
		).
		Values(
			media.PostID,
			media.Type,
			media.URL,
			media.ThumbnailURL,
		).
		Suffix(fmt.Sprintf("RETURNING %s", mediaIdColumn))

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error(op, "error in build sql media", "error", err)
		return uuid.Nil, fmt.Errorf(op+" error in create media %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var id uuid.UUID

	if err := r.db.DB().ScanOneContext(ctx, &id, q, args...); err != nil {
		logger.Error(op, "error in create media", "error", err)
		return uuid.Nil, fmt.Errorf(op+" error in create media %w", err)
	}

	return id, nil
}
