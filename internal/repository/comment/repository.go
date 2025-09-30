package comment

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx"

	"post/internal/client/db"
	"post/internal/model"
	"post/internal/repository"
)

const (
	commentIdColumn       = "comment_id"
	postIdColumn          = "post_id"
	userIdColumn          = "user_id"
	parentCommentIdColumn = "parent_comment_id"
	replyToUserIdColumn   = "reply_to_user_id"
	contentColumn         = "content"
	createdAtColumn       = "created_at"
	updatedAtColumn       = "updated_at"
)

const (
	commentsTable = "comments"
)

// Модели репозитория для внутренних операций
type Comment struct {
	ID              uuid.UUID  `db:"comment_id"`
	PostID          uuid.UUID  `db:"post_id"`
	UserID          uuid.UUID  `db:"user_id"`
	ParentCommentID *uuid.UUID `db:"parent_comment_id"`
	ReplyToUserID   *uuid.UUID `db:"reply_to_user_id"`
	Content         string     `db:"content"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
}

type repo struct {
	db db.Client
}

// Add implements repository.CommentRepository - создание комментария или ответа
func (r repo) Add(ctx context.Context, postID, userID uuid.UUID, content string, parentCommentID *uuid.UUID) (*model.Comment, error) {
	const op = "repository.comment.Add"

	// Если это ответ на комментарий, нужно получить ID пользователя, которому отвечают
	var replyToUserID *uuid.UUID
	if parentCommentID != nil {
		parentComment, err := r.getCommentByID(ctx, *parentCommentID)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to get parent comment: %w", op, err)
		}
		replyToUserID = &parentComment.UserID
	}

	builder := sq.Insert(commentsTable).
		Columns(postIdColumn, userIdColumn, contentColumn, parentCommentIdColumn, replyToUserIdColumn).
		Values(postID, userID, content, parentCommentID, replyToUserID).
		Suffix("RETURNING " + commentIdColumn + ", " + createdAtColumn + ", " + updatedAtColumn).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var result struct {
		CommentID uuid.UUID `db:"comment_id"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}

	if err := r.db.DB().ScanOneContext(ctx, &result, q, args...); err != nil {
		return nil, fmt.Errorf("%s: failed to create comment: %w", op, err)
	}

	// Возвращаем модель в формате API
	return &model.Comment{
		ID:              result.CommentID,
		PostID:          postID,
		UserID:          userID,
		ParentCommentID: parentCommentID,
		ReplyToUserID:   replyToUserID,
		Content:         content,
		CreatedAt:       result.CreatedAt,
		UpdatedAt:       result.UpdatedAt,
		IsReply:         parentCommentID != nil,
	}, nil
}

// getCommentByID - вспомогательный метод для получения комментария по ID
func (r repo) getCommentByID(ctx context.Context, commentID uuid.UUID) (*Comment, error) {
	const op = "repository.comment.getCommentByID"

	builder := sq.Select("*").
		From(commentsTable).
		Where(sq.Eq{commentIdColumn: commentID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	q := db.Query{
		Name:     op,
		QueryRaw: query,
	}

	var comment Comment
	if err := r.db.DB().ScanOneContext(ctx, &comment, q, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: comment not found", op)
		}
		return nil, fmt.Errorf("%s: failed to get comment: %w", op, err)
	}

	return &comment, nil
}

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

// convertToModelComment - конвертация внутренней модели в API модель (без данных пользователя)
func (r repo) convertToModelComment(comment *Comment) *model.Comment {
	return &model.Comment{
		ID:              comment.ID,
		PostID:          comment.PostID,
		UserID:          comment.UserID,
		ParentCommentID: comment.ParentCommentID,
		ReplyToUserID:   comment.ReplyToUserID,
		Content:         comment.Content,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		User:            nil, // Будет заполнено в сервисе
		Replies:         nil, // Будет заполнено позже
		IsReply:         comment.ParentCommentID != nil,
	}
}

// GetReplies implements repository.CommentRepository - получение ответов на конкретный комментарий
func (r repo) GetReplies(ctx context.Context, commentID uuid.UUID) ([]*model.Comment, error) {
	const op = "repository.comment.GetReplies"

	// Получаем ответы на комментарий
	repliesQuery := sq.Select("*").
		From(commentsTable).
		Where(sq.Eq{parentCommentIdColumn: commentID}).
		OrderBy(createdAtColumn + " ASC").
		PlaceholderFormat(sq.Dollar)

	repliesSql, repliesArgs, err := repliesQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to build replies query: %w", op, err)
	}

	repliesQ := db.Query{
		Name:     op,
		QueryRaw: repliesSql,
	}

	var replies []Comment
	if err := r.db.DB().ScanAllContext(ctx, &replies, repliesQ, repliesArgs...); err != nil {
		return nil, fmt.Errorf("%s: failed to get replies: %w", op, err)
	}

	// Конвертируем в API модель
	var result []*model.Comment
	for _, reply := range replies {
		result = append(result, r.convertToModelComment(&reply))
	}

	return result, nil
}

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

func New(db db.Client) repository.CommentRepository {
	return &repo{db: db}
}
