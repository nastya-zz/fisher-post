package comment

import (
	"context"
	"github.com/google/uuid"
	"post/internal/model"
	"post/internal/repository"
)
import "post/internal/service"

type serv struct {
	commentRepository repository.CommentRepository
}

func (s serv) AddComment(ctx context.Context, postID, userID uuid.UUID) (*model.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s serv) RemoveComment(ctx context.Context, postID, userID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func New(commentRepository repository.CommentRepository) service.CommentService {
	return &serv{commentRepository}
}
