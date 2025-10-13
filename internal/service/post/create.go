package post

import (
	"context"
	"encoding/json"
	"fmt"

	"post/internal/model"
	repoModel "post/internal/repository/post/model"
	"post/pkg/logger"
)

func (s serv) CreatePost(ctx context.Context, post *model.CreatePost) (*model.Post, error) {
	const op = "service.post.CreatePost"

	var createdPost *repoModel.CreatedPost
	var postModel *model.Post
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		createdPost, errTx = s.repository.CreatePost(ctx, post)
		if errTx != nil {
			return fmt.Errorf("%s %w", op, errTx)
		}

		for _, fishId := range post.FishTypeIDs {
			errTx = s.repository.CreatePostFishReference(ctx, createdPost.ID, fishId)
		}
		if errTx != nil {
			logger.Error("failed to create post fish reference", "error", errTx)
			return fmt.Errorf("%s %w", op, errTx)
		}

		for _, tackleId := range post.TackleTypeIDs {
			errTx = s.repository.CreatePostTackleReference(ctx, createdPost.ID, tackleId)
		}
		if errTx != nil {
			logger.Error("failed to create post tackle reference", "error", errTx)
			return fmt.Errorf("%s %w", op, errTx)
		}

		user, err := s.userService.GetUser(ctx, "", createdPost.UserID)
		if err != nil {
			logger.Error("failed to get user", "error", err)
			return fmt.Errorf("%s %w", op, fmt.Errorf("failed to get user: %w", err))
		}

		postModel = &model.Post{
			ID: createdPost.ID,
			User: model.User{
				ID:        user.ID,
				Username:  user.Username,
				AvatarUrl: user.AvatarUrl,
			},
			Description: createdPost.Description,
			Geolocation: model.Geolocation{
				Latitude:  createdPost.Latitude,
				Longitude: createdPost.Longitude,
			},
			CreatedAt: createdPost.CreatedAt,
		}

		body, errTx := json.Marshal(model.PostPayload{
			ID:       createdPost.ID.String(),
			AuthorID: createdPost.UserID.String(),
		})

		if errTx != nil {
			return fmt.Errorf("error in marshal json body %w", errTx)
		}

		event := &model.Event{
			Type:    model.PostCreated,
			Payload: body,
		}

		errTx = s.eventRepository.SaveEvent(ctx, event)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	return postModel, err
}
