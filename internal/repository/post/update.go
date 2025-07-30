package post

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"post/internal/client/db"
	"post/internal/model"
)

func (r repo) Update(ctx context.Context, post *model.UpdatePost) error {
	const op = "post.repository.update"

	query := sq.Update(postsTable).
		Set(descriptionColumn, post.Description).
		Set(latitudeColumn, post.Geolocation.Latitude).
		Set(longitudeColumn, post.Geolocation.Longitude).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{postIdColumn: post.ID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf(op+" failed to build query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: sql,
	}

	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf(op+" failed to update post: %w", err)
	}

	return nil

}
