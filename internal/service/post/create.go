package post

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	repoModel "post/internal/repository/post/model"
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
			return fmt.Errorf("%s %w", op, errTx)
		}

		for _, tackleId := range post.TackleTypeIDs {
			errTx = s.repository.CreatePostTackleReference(ctx, createdPost.ID, tackleId)
		}
		if errTx != nil {
			return fmt.Errorf("%s %w", op, errTx)
		}

		userResponse, err := s.userService.GetUser(ctx, "", createdPost.UserID)
		if err != nil {
			return fmt.Errorf("%s %w", op, fmt.Errorf("failed to get user: %w", err))
		}

		user := userResponse.GetProfile()

		postModel = &model.Post{
			ID: createdPost.ID,
			User: model.User{
				ID:        uuid.MustParse(user.Id),
				Username:  user.GetName(),
				AvatarUrl: user.GetAvatarPath(),
			},
			Description: createdPost.Description,
			Geolocation: model.Geolocation{
				Latitude:  createdPost.Latitude,
				Longitude: createdPost.Longitude,
			},
			CreatedAt: createdPost.CreatedAt,
		}

		return nil
	})

	return postModel, err
}
