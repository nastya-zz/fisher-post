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


func (UnimplementedPostServiceServer) AddComment(context.Context, *AddCommentRequest) (*Comment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddComment not implemented")
}
func (desc.UnimplementedPostServiceServer) GetComments(context.Context, *desc.PostId) (*desc.GetCommentsResponse, error)

func (UnimplementedPostServiceServer) RemoveComment(context.Context, *RemoveCommentRequest) (*RemoveCommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveComment not implemented")
}
*/
