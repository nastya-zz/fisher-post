package post

import (
	"context"

	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/converter"
	"post/internal/model"
	"post/pkg/helpers"
)

func (i *Implementation) UpdatePost(ctx context.Context, req *desc.UpdatePostRequest) (*desc.Post, error) {

	postID, err := model.GetUuid(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	updatePost := &model.UpdatePost{
		ID:          postID,
		Description: req.Description,
		Geolocation: model.Geolocation{
			Latitude:  req.Location.Latitude,
			Longitude: req.Location.Longitude,
		},
		FishTypeIDs:   helpers.GetIntList(req.FishTypeIds),
		TackleTypeIDs: helpers.GetIntList(req.TackleTypeIds),
	}

	post, err := i.postService.UpdatePost(ctx, updatePost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", "Ошибка при обновлении поста")
	}

	return converter.FromModelPostToDescPost(post), nil
}
