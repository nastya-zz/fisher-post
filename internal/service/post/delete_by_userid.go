package post

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/pkg/logger"
)

func (s serv) DeleteByUserId(ctx context.Context, userId uuid.UUID) error {
	const op = "service.post.DeleteByUserId"

	s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		posts, err := s.repository.GetPosts(ctx, userId)
		if err != nil {
			return fmt.Errorf(op+" failed to get posts: %w", err)
		}
		logger.Info(op, "user id", userId, "posts", posts)

		for _, post := range posts {
			media, err := s.mediaRepository.GetByPostID(ctx, post.ID)
			if err != nil {
				return fmt.Errorf(op+" failed to get media: %w", err)
			}
			logger.Info(op, "media", media)

			var link []string
			for _, m := range media {
				link = append(link, m.Url)
			}

			var errors []string
			for _, l := range link {
				err = s.minioService.RemoveFile(ctx, l)
				if err != nil {
					errors = append(errors, fmt.Sprintf("failed to remove file: %s", err.Error()))
				}
			}
			if len(errors) > 0 {
				logger.Error(op, "failed to remove files", "errors", errors)
			}
		}

		err = s.repository.DeleteByUserId(ctx, userId)
		if err != nil {
			return fmt.Errorf(op+" failed to delete post: %w", err)
		}

		return nil
	})
	return nil
}
