package media

import (
	"post/internal/client/db"
	"post/internal/repository"
)

type repo struct {
	db db.Client
}

const (
	mediaTable              = "media"
	mediaIdColumn           = "media_id"
	postIdColumn            = "post_id"
	mediaTypeColumn         = "media_type"
	mediaUrlColumn          = "url"
	mediaThumbnailUrlColumn = "thumbnail_url"
	mediaCreatedAtColumn    = "created_at"
)

func New(db db.Client) repository.MediaRepository {
	return &repo{db: db}
}
