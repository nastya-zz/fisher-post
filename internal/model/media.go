package model

import (
	"time"

	"github.com/google/uuid"
)

type DescCreateMedia struct {
	File        []byte
	PostID      uuid.UUID
	Filename    string
	Type        string
	IsThumbnail bool
}

type CreateMedia struct {
	PostID uuid.UUID
	MediaID uuid.UUID
	URL string
	ThumbnailURL string
	Type        string
	IsThumbnail bool
	CreatedAt time.Time
}