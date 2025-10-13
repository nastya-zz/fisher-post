package post

import (
	"context"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/model"
	"post/internal/utils"
)

func (i *Implementation) DeletePost(ctx context.Context, req *desc.DeletePostRequest) (*desc.DeletePostResponse, error) {

	userID, err := utils.GetUserIdFromMetadata(ctx)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id пользователя не валидный")
	}

	postID, err := model.GetUuid(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	if postID == uuid.Nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	err = i.postService.DeletePost(ctx, postID, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", "Ошибка при удалении поста")
	}

	return &desc.DeletePostResponse{
		Success: true,
	}, nil
}
