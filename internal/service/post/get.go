package post

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"post/internal/model"
	"post/pkg/logger"
)

func (s serv) GetPost(ctx context.Context, id uuid.UUID) (*model.Post, error) {
	const op = "service.post.GetPost"

	post, err := s.repository.Get(ctx, id)

	if err != nil {
		logger.Error(op+"failed to get post", "err", err)
		return nil, fmt.Errorf("failed to get post: %w", err)
	}

	if post == nil {
		return nil, fmt.Errorf("post not found")
	}

	userResponse, err := s.userService.GetUser(ctx, "", post.UserID)
	if err != nil {
		logger.Error("failed to get user", "error", err)
		return nil, fmt.Errorf("%s %w", op, fmt.Errorf("failed to get user: %w", err))
	}

	user := userResponse.GetProfile()

	likesCount, err := s.likeRepository.GetLikesCount(ctx, post.ID)
	if err != nil {
		logger.Error("failed to get likes count", "error", err)
		return nil, fmt.Errorf("%s %w", op, fmt.Errorf("failed to get likes count: %w", err))
	}

	postModel := &model.Post{
		ID: post.ID,
		User: model.User{
			ID:        uuid.MustParse(user.Id),
			Username:  user.GetName(),
			AvatarUrl: user.GetAvatarPath(),
		},
		Description: post.Description,
		Geolocation: model.Geolocation{
			Latitude:  post.Latitude,
			Longitude: post.Longitude,
		},
		TackleTypes: func() []model.Dictionary {
			result := make([]model.Dictionary, 0, len(post.TackleTypes))
			for _, tackle := range post.TackleTypes {
				result = append(result, model.Dictionary{
					ID:          tackle.ID,
					Name:        tackle.Name,
					Description: tackle.Description,
				})
			}
			return result
		}(),
		FishTypes: func() []model.Dictionary {
			result := make([]model.Dictionary, 0, len(post.FishTypes))
			for _, fish := range post.FishTypes {
				result = append(result, model.Dictionary{
					ID:          fish.ID,
					Name:        fish.Name,
					Description: fish.Description,
				})
			}
			return result
		}(),
		CreatedAt: post.CreatedAt,
		LikesCount: likesCount,
	}

	return postModel, nil
}
