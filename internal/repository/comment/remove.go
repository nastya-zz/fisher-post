package comment

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"post/internal/client/db"
)

// Remove implements repository.CommentRepository - удаление комментария
func (r repo) Remove(ctx context.Context, commentID, userID uuid.UUID) error {
	const op = "repository.comment.Remove"

	// Проверяем, что комментарий существует и принадлежит пользователю
	comment, err := r.getCommentByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("%s: failed to get comment: %w", op, err)
	}

	// Проверяем права доступа - только автор может удалить свой комментарий
	if comment.UserID != userID {
		return fmt.Errorf("%s: access denied - comment belongs to another user", op)
	}

	// Удаляем комментарий
	builder := sq.Delete(commentsTable).
		Where(sq.Eq{commentIdColumn: commentID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	if _, err := r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("%s: failed to delete comment: %w", op, err)
	}

	return nil
}
