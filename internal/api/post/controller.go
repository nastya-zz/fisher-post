package post

import (
	desc "github.com/nastya-zz/fisher-protocols/gen/post_v1"
	"post/internal/service"
)

type Implementation struct {
	desc.UnimplementedPostServiceServer
	postService    service.PostService
	commentService service.CommentService
}

func NewImplementation(postService service.PostService, commentService service.CommentService) *Implementation {
	return &Implementation{
		postService:    postService,
		commentService: commentService,
	}
}
