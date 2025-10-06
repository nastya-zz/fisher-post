package comment

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	"post/pkg/logger"
)

// enrichCommentsWithUsers - обогащает комментарии информацией о пользователях
// Возвращает только комментарии, для которых удалось получить данные пользователя
func (s serv) enrichCommentsWithUsers(ctx context.Context, comments []*model.Comment) ([]*model.Comment, error) {
	var enrichedComments []*model.Comment

	for _, comment := range comments {
		// Пропускаем комментарии без указанного пользователя
		if comment.UserID == uuid.Nil {
			logger.Debug("skipping comment with nil userID", "commentID", comment.ID)
			continue
		}

		// Получаем данные пользователя
		user, err := s.userService.GetUser(ctx, "", comment.UserID)
		if err != nil {
			logger.Error("failed to get user for comment", "error", err, "userID", comment.UserID)
			// Пропускаем комментарий, если не удалось получить данные пользователя
			continue
		}

		comment.User = user

		// Рекурсивно обогащаем ответы
		if len(comment.Replies) > 0 {
			enrichedReplies, err := s.enrichCommentsWithUsers(ctx, comment.Replies)
			if err != nil {
				return nil, fmt.Errorf("failed to enrich replies for comment %s: %w", comment.ID, err)
			}
			comment.Replies = enrichedReplies
		}

		// Добавляем успешно обогащенный комментарий в результат
		enrichedComments = append(enrichedComments, comment)
	}

	return enrichedComments, nil
}
