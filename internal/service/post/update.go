package post

import (
	"context"
	"fmt"

	"post/internal/model"
	"post/pkg/logger"
)

func (s serv) UpdatePost(ctx context.Context, post *model.UpdatePost) (*model.Post, error) {
	const op = "post.service.update"

	//todo: check if user is the owner of the post

	existPost, err := s.repository.Get(ctx, post.ID)
	if err != nil && existPost == nil {
		logger.Error(op+" post not found", "error", err)
		return nil, fmt.Errorf(op+" post not found: %w", err)
	}

	var postModel *model.Post

	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		errTx = s.repository.Update(ctx, post)
		if errTx != nil {
			return fmt.Errorf(op+" failed to update post: %w", errTx)
		}

		for _, fishId := range post.FishTypeIDs {
			errTx = s.repository.UpdatePostFishReference(ctx, post.ID, fishId)
			if errTx != nil {
				return fmt.Errorf(op+" failed to update post fish reference: %w", errTx)
			}
		}

		for _, tackleId := range post.TackleTypeIDs {
			errTx = s.repository.UpdatePostTackleReference(ctx, post.ID, tackleId)
			if errTx != nil {
				return fmt.Errorf(op+" failed to update post tackle reference: %w", errTx)
			}
		}

		postModel, errTx = s.GetPost(ctx, post.ID)
		if errTx != nil {
			return fmt.Errorf(op+" failed to get post: %w", errTx)
		}

		return nil
	})

	if err != nil {
		logger.Error(op+" failed to update post", "error", err)
		return nil, fmt.Errorf(op+" failed to update post: %w", err)
	}

	return postModel, nil
}
