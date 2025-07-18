package model

import (
	"time"

	"github.com/google/uuid"
)

type CreatedPost struct {
	ID          uuid.UUID `db:"post_id"`
	UserID      uuid.UUID `db:"user_id"`
	Description string    `db:"description"`
	Latitude    float64   `db:"latitude"`
	Longitude   float64   `db:"longitude"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Post struct {
	ID          uuid.UUID    `db:"post_id"`
	UserID      uuid.UUID    `db:"user_id"`
	Description string       `db:"description"`
	Latitude    float64      `db:"latitude"`
	Longitude   float64      `db:"longitude"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	TackleTypes []TackleType `db:"tackle_types"`
	FishTypes   []FishType   `db:"fish_types"`
}

type TackleType struct {
	ID          int    `db:"tackle_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type FishType struct {
	ID          int    `db:"fish_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
