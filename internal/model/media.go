package model

import "github.com/google/uuid"

type CreateMedia struct {
	File []byte
	UserID uuid.UUID
	Filename string
	Type string
	IsThumbnail bool
}