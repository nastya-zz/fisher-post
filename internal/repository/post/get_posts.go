package post

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
	repoModel "post/internal/repository/post/model"
)

func (r repo) GetPosts(ctx context.Context, id uuid.UUID) ([]*repoModel.Post, error) {
	const op = "repository.post.GetPosts"

	query := sq.Select("*").From(postsTable).Where(sq.Eq{userIdColumn: id}).PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf(op+" failed to build query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: sql,
	}

	var posts []*repoModel.Post
	if err := r.db.DB().ScanAllContext(ctx, &posts, q, args...); err != nil {
		return nil, fmt.Errorf(op+" failed to get posts: %w", err)
	}

	return posts, nil	
}
