package post

import (
	"context"

	"github.com/google/uuid"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) GetLikes(ctx context.Context, req *desc.PostId) (*desc.GetLikesResponse, error) {
	const op = "api.post.GetLikes"

	if req.GetPostId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "post id is required")
	}

	postID, err := uuid.Parse(req.GetPostId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post id: %v", err)
	}

	likes, err := i.likeService.GetLikes(ctx, postID)
	if err != nil {
		return nil, err
	}

	var users []*desc.User
	for _, like := range likes {
		users = append(users, &desc.User{
			Id:        like.ID.String(),
			Username:  like.Username,
			AvatarUrl: like.AvatarUrl,
		})
	}

	return &desc.GetLikesResponse{
		Likes: users,
	}, nil
}
