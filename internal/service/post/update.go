package post

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	"post/pkg/logger"
)

func (s serv) UpdatePost(ctx context.Context, post *model.UpdatePost, userID uuid.UUID) (*model.Post, error) {
	const op = "post.service.update"

	existPost, err := s.repository.Get(ctx, post.ID)
	if err != nil && existPost == nil {
		logger.Error(op+" post not found", "error", err)
		return nil, fmt.Errorf(op+" post not found: %w", err)
	}

	if existPost.UserID != userID {
		logger.Error(op+" user is not the owner of the post", "error", err)
		return nil, fmt.Errorf(op+" user is not the owner of the post: %w", err)
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

		body, errTx := json.Marshal(model.PostPayload{
			ID:       postModel.ID.String(),
			AuthorID: postModel.User.ID.String(),
		})

		if errTx != nil {
			return fmt.Errorf("error in marshal json body %w", errTx)
		}

		event := &model.Event{
			Type:    model.PostUpdated,
			Payload: body,
		}

		errTx = s.eventRepository.SaveEvent(ctx, event)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		logger.Error(op+" failed to update post", "error", err)
		return nil, fmt.Errorf(op+" failed to update post: %w", err)
	}

	return postModel, nil
}
