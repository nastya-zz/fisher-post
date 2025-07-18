package post

import (
	"context"

	"github.com/google/uuid"

	"post/internal/client/db"
	"post/internal/model"
	"post/internal/repository"
)

const (
	postIdColumn = "post_id"
)

const (
	postsTable        = "posts"
	userIdColumn      = "user_id"
	descriptionColumn = "description"
	latitudeColumn    = "latitude"
	longitudeColumn   = "longitude"
	createdAtColumn   = "created_at"
	updatedAtColumn   = "updated_at"
)

const (
	postFishTable = "post_fish"
	fishIdColumn  = "fish_id"
)

const (
	postTackleTable = "post_tackle"
	tackleIdColumn  = "tackle_id"
)

const (
	mediaTable           = "media"
	mediaIdColumn        = "media_id"
	mediaUrlColumn       = "url"
	thumbnailUrlColumn   = "thumbnail_url"
	mediaTypeColumn      = "media_type"
	mediaSizeColumn      = "size"
	mediaCreatedAtColumn = "created_at"
)

const (
	tacleTypesTable = "tackle_types"
	fishTypesTable  = "fish_types"
)

const (
	dictionaryNameColumn        = "name"
	dictionaryDescriptionColumn = "description"
)

type repo struct {
	db db.Client
}

func (r repo) Update(ctx context.Context, post *model.Post) (*model.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (r repo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}


func New(db db.Client) repository.PostRepository {
	return &repo{db: db}
}
