package model

import (
	"github.com/google/uuid"
	"time"
)

type user struct {
	ID        uuid.UUID
	Username  string
	AvatarUrl string
}

const MediaTypePhoto = "PHOTO"
const MediaTypeVideo = "VIDEO"

type media struct {
	ID           uuid.UUID
	MediaType    string
	Url          string
	ThumbnailUrl string
}

type Geolocation struct {
	Latitude  float64
	Longitude float64
}

type Post struct {
	ID          uuid.UUID
	User        user
	Description string
	geolocation Geolocation
	CreatedAt   time.Time
	Media       []media
	Likes       int
	Comments    []Comment
	FishType    []Dictionary
	TackleType  []Dictionary
}

func GetUuid[T ~string](id T) (uuid.UUID, error) {
	return uuid.Parse(string(id))
}
