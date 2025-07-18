package like

import (
	"post/internal/client/db"
	"post/internal/repository"
)

const (
	likesTable = "likes"
)

const (
	postIdColumn    = "post_id"
	userIdColumn    = "user_id"
	createdAtColumn = "created_at"
	likeIdColumn    = "like_id"
)

type repo struct {
	db db.Client
}

func New(db db.Client) repository.LikeRepository {
	return &repo{db: db}
}
