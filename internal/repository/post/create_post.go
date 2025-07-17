package post

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"

	"post/internal/client/db"
	"post/internal/model"
	repoModel "post/internal/repository/post/model"
	"post/pkg/logger"
)

func (r repo) CreatePost(ctx context.Context, post *model.CreatePost) (*repoModel.CreatedPost, error) {
	const op = "repository.post.Create"

	builder := sq.Insert(postsTable).
		Columns(userIdColumn, descriptionColumn, latitudeColumn, longitudeColumn).
		PlaceholderFormat(sq.Dollar).
		Values(post.UserID, post.Description, post.Geolocation.Latitude, post.Geolocation.Longitude).
		Suffix("RETURNING *")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var createdPost repoModel.CreatedPost

	if err := r.db.DB().ScanOneContext(ctx, &createdPost, q, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(op+" error in create post %w", err)
		}

		logger.Warn(fmt.Errorf("error in create post %w", err).Error())
		return nil, fmt.Errorf(op+" error in create post %w", err)
	}

	return &createdPost, nil

}
