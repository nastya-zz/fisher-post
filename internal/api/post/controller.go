package post

import (
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"

	"post/internal/service"
)

type Implementation struct {
	desc.UnimplementedPostServiceServer
	postService    service.PostService
	commentService service.CommentService
	likeService    service.LikeService
}

func NewImplementation(postService service.PostService, commentService service.CommentService, likeService service.LikeService) *Implementation {
	return &Implementation{
		postService:    postService,
		commentService: commentService,
		likeService:    likeService,
	}
}

/*
func (UnimplementedPostServiceServer) CreatePost(context.Context, *CreatePostRequest) (*Post, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePost not implemented")
}

func (UnimplementedPostServiceServer) UpdatePost(context.Context, *UpdatePostRequest) (*Post, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePost not implemented")
}
func (UnimplementedPostServiceServer) DeletePost(context.Context, *DeletePostRequest) (*DeletePostResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePost not implemented")
}
func (UnimplementedPostServiceServer) AddLike(context.Context, *AddLikeRequest) (*Like, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLike not implemented")
}
func (UnimplementedPostServiceServer) AddComment(context.Context, *AddCommentRequest) (*Comment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddComment not implemented")
}
func (UnimplementedPostServiceServer) RemoveLike(context.Context, *RemoveLikeRequest) (*RemoveLikeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveLike not implemented")
}
func (UnimplementedPostServiceServer) RemoveComment(context.Context, *RemoveCommentRequest) (*RemoveCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveComment not implemented")
}
*/
