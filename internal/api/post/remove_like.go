package post

import (
	"context"

	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/model"
)

func (i *Implementation) RemoveLike(ctx context.Context, req *desc.RemoveLikeRequest) (*desc.RemoveLikeResponse, error) {
	userID, err := model.GetUuid(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id пользователя не валидный")
	}

	postID, err := model.GetUuid(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	likesCount, err := i.likeService.RemoveLike(ctx, postID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &desc.RemoveLikeResponse{
		Success:       true,
		NewLikesCount: int32(likesCount),
	}, nil
}
