package post

import (
	"context"

	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/model"
	"post/internal/utils"
)

func (i *Implementation) AddLike(ctx context.Context, req *desc.AddLikeRequest) (*desc.AddLikeResponse, error) {

	userID, err := utils.GetUserIdFromMetadata(ctx)
	if err != nil {
		return nil, err
	}

	postID, err := model.GetUuid(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	likesCount, err := i.likeService.AddLike(ctx, postID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &desc.AddLikeResponse{
		Success:    true,
		LikesCount: int32(likesCount),
	}, nil
}
