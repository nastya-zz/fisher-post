package model

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID        uuid.UUID
	Content   string
	CreatedAt time.Time
}
