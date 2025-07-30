package post

import (
	"context"

	"github.com/nastya-zz/fisher-protocols/gen/post_v1"
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"post/internal/converter"
	"post/internal/model"
)

func (i *Implementation) GetPosts(ctx context.Context, req *desc.GetPostRequest) (*desc.GetPostsResponse, error) {
	postID, err := model.GetUuid(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", "Id поста не валидный")
	}

	posts, err := i.postService.GetPosts(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", "Ошибка при получении постов")
	}

	pbPosts := make([]*post_v1.Post, len(posts))
	for i, p := range posts {
		pbPosts[i] = converter.FromModelPostToDescPost(p)
	}

	return &desc.GetPostsResponse{
		Posts: pbPosts,
	}, nil
}
