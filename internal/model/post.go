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
	ID           uuid.UUID
	MediaType    string
	Url          string
	ThumbnailUrl string
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

func GetUuid[T ~string](id T) (uuid.UUID, error) {
	return uuid.Parse(string(id))
}
