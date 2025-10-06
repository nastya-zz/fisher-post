package comment

import (
	"time"

	"github.com/google/uuid"

	"post/internal/client/db"
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

// Comment - модель репозитория для внутренних операций
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

func New(db db.Client) repository.CommentRepository {
	return &repo{db: db}
}
