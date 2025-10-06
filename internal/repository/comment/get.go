package comment

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
)

// GetByPostID implements repository.CommentRepository - получение всех комментариев к посту с вложенностью
func (r repo) GetByPostID(ctx context.Context, postID uuid.UUID, limit, offset int) ([]*model.Comment, error) {
	const op = "repository.comment.GetByPostID"

	// Валидация параметров пагинации
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	// Получаем корневые комментарии (без родителя)
	rootCommentsQuery := sq.Select("*").
		From(commentsTable).
		Where(sq.Eq{postIdColumn: postID}).
		Where(sq.Eq{parentCommentIdColumn: nil}).
		OrderBy(createdAtColumn + " ASC").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar)

	rootQuery, rootArgs, err := rootCommentsQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build root comments query: %w", op, err)
	}

	rootQ := db.Query{
		Name:     op + ".GetRootComments",
		QueryRaw: rootQuery,
	}

	var rootComments []Comment
	if err := r.db.DB().ScanAllContext(ctx, &rootComments, rootQ, rootArgs...); err != nil {
		return nil, fmt.Errorf("%s: failed to get root comments: %w", op, err)
	}

	// Для каждого корневого комментария получаем ответы
	var result []*model.Comment
	for _, rootComment := range rootComments {
		comment := r.convertToModelComment(&rootComment)

		// Получаем ответы на этот комментарий
		replies, err := r.GetReplies(ctx, rootComment.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get replies for comment %s: %w", op, rootComment.ID, err)
		}

		comment.Replies = replies
		result = append(result, comment)
	}

	return result, nil
}

// GetCommentsCount implements repository.CommentRepository - подсчет комментариев
func (r repo) GetCommentsCount(ctx context.Context, postID uuid.UUID) (int, error) {
	const op = "repository.comment.GetCommentsCount"

	builder := sq.Select("COUNT(*)").
		From(commentsTable).
		Where(sq.Eq{postIdColumn: postID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var count int
	if err := r.db.DB().ScanOneContext(ctx, &count, q, args...); err != nil {
		return 0, fmt.Errorf("%s: failed to get comments count: %w", op, err)
	}

	return count, nil
}
