package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID  `db:"comment_id" json:"id"`
	PostID          uuid.UUID  `db:"post_id" json:"post_id"`
	UserID          uuid.UUID  `db:"user_id" json:"user_id"`
	ParentCommentID *uuid.UUID `db:"parent_comment_id" json:"parent_comment_id"` // NULL для корневых комментариев
	ReplyToUserID   *uuid.UUID `db:"reply_to_user_id" json:"reply_to_user_id"`   // ID пользователя, которому отвечают
	Content         string     `db:"content" json:"content"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`

	// Для API ответов добавим связанные данные
	User    *User      `json:"user,omitempty"`    // Автор комментария
	Replies []*Comment `json:"replies,omitempty"` // Ответы на этот комментарий (если есть)
	IsReply bool       `json:"is_reply"`          // Является ли ответом на другой комментарий
}
