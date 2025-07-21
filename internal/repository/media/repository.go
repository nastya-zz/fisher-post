package media

import (
	"context"

	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
	"post/internal/repository"
)

type repo struct {
	db db.Client
}

// Create implements repository.MediaRepository.
func (r *repo) Create(ctx context.Context, media *model.CreateMedia) (uuid.UUID, error) {
	panic("unimplemented")
}

// Delete implements repository.MediaRepository.
func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// Get implements repository.MediaRepository.
func (r *repo) Get(ctx context.Context, id uuid.UUID) (*model.Media, error) {
	panic("unimplemented")
}

// GetByPostID implements repository.MediaRepository.
func (r *repo) GetByPostID(ctx context.Context, postID uuid.UUID) ([]model.Media, error) {
	panic("unimplemented")
}

func New(db db.Client) repository.MediaRepository {
	return &repo{db: db}
}
