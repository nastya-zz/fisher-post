package post

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/pkg/logger"
)

func (s serv) DeletePost(ctx context.Context, id uuid.UUID) error {
	const op = "service.post.DeletePost"

	media, err := s.mediaRepository.GetByPostID(ctx, id)
	if err != nil {
		return fmt.Errorf(op+" failed to get media: %w", err)
	}
	logger.Info(op, "media", media)

	var link []string
	for _, m := range media {
		link = append(link, m.Url)
	}

	err = s.repository.Delete(ctx, id)
	if err != nil {
		logger.Error(op, "failed to delete post", "error", err)
		return fmt.Errorf(op+" failed to delete post: %w", err)
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

	return nil
}
