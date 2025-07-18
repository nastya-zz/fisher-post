package post

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx"

	"post/internal/client/db"
	repoModel "post/internal/repository/post/model"
)

func (r repo) Get(ctx context.Context, id uuid.UUID) (*repoModel.Post, error) {
	const op = "repository.post.Get"

	builder := sq.Select("*").
		From(postsTable).
		Where(sq.Eq{postIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf(op+" failed to build query: %w", err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var post repoModel.Post
	if err := r.db.DB().ScanOneContext(ctx, &post, q, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf(op+" failed to get post: %w", err)
	}

	// Получаем связанные типы рыб
	fishBuilder := sq.Select("d.fish_id", "d.name", "d.description").
		From(fishTypesTable + " d").
		Join(postFishTable + " pf ON pf.fish_id = d.fish_id").
		Where(sq.Eq{"pf." + postIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	fishQuery, fishArgs, err := fishBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf(op+" failed to build fish query: %w", err)
	}

	q = db.Query{
		Name:     op + ".GetFishTypes",
		QueryRaw: fishQuery,
	}

	if err := r.db.DB().ScanAllContext(ctx, &post.FishTypes, q, fishArgs...); err != nil {
		return nil, fmt.Errorf(op+" failed to get fish types: %w", err)
	}

	// Получаем связанные типы снастей
	tackleBuilder := sq.Select("d.tackle_id", "d.name", "d.description").
		From(tacleTypesTable + " d").
		Join(postTackleTable + " pf ON pf.tackle_id = d.tackle_id").
		Where(sq.Eq{"pf." + postIdColumn: id}).
		PlaceholderFormat(sq.Dollar)

	tackleQuery, tackleArgs, err := tackleBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf(op+" failed to build tackle query: %w", err)
	}

	q = db.Query{
		Name:     op + ".GetTackleTypes",
		QueryRaw: tackleQuery,
	}

	if err := r.db.DB().ScanAllContext(ctx, &post.TackleTypes, q, tackleArgs...); err != nil {
		return nil, fmt.Errorf(op+" failed to get tackle types: %w", err)
	}

	return &post, nil
}
