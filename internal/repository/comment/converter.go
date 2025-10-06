package comment

import (
	"post/internal/model"
)

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
