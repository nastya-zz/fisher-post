package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Username  string
	AvatarUrl string
}

const MediaTypePhoto = "PHOTO"
const MediaTypeVideo = "VIDEO"

type Media struct {
	ID           uuid.UUID `db:"media_id"`
	MediaType    string    `db:"media_type"`
	Url          string    `db:"url"`
	ThumbnailUrl string    `db:"thumbnail_url"`
}

type Geolocation struct {
	Latitude  float64
	Longitude float64
}

type CreatePost struct {
	UserID        uuid.UUID
	Description   string
	Geolocation   Geolocation
	FishTypeIDs   []int
	TackleTypeIDs []int
}

type Post struct {
	ID            uuid.UUID
	User          User
	Description   string
	Geolocation   Geolocation
	CreatedAt     time.Time
	Media         []Media
	LikesCount    int
	CommentsCount int
	FishTypes     []Dictionary
	TackleTypes   []Dictionary
}

type UpdatePost struct {
	ID            uuid.UUID
	Description   string
	Geolocation   Geolocation
	FishTypeIDs   []int
	TackleTypeIDs []int
}

func GetUuid[T ~string](id T) (uuid.UUID, error) {
	return uuid.Parse(string(id))
}
