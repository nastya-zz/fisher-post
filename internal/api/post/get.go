package post

import (
	"context"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"post/internal/converter"
	"post/internal/model"
)

func (i *Implementation) GetPost(ctx context.Context, req *desc.GetPostRequest) (*desc.Post, error) {
	postID, err := model.GetUuid(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	if postID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	post, err := i.postService.GetPost(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", "Ошибка при получении поста")
	}

	return &desc.Post{
		Id:          post.ID.String(),
		User:        converter.FromModelUserToDescUser(post.User),
		Description: post.Description,
		Location: &desc.LatLng{
			Latitude:  post.Geolocation.Latitude,
			Longitude: post.Geolocation.Longitude,
		},
		Media:         converter.FromModelMediaToDescMedia(post.Media),
		LikesCount:    int32(post.LikesCount),
		CommentsCount: int32(post.CommentsCount),
		FishTypes:     converter.FromFishTypesToDescFishTypes(post.FishTypes),
		TackleTypes:   converter.FromFishTypesToDescTackleType(post.TackleTypes),
		CreatedAt:     timestamppb.New(post.CreatedAt),
	}, nil
}
