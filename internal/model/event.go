package model

import "time"

type Event struct {
	ID      int
	Type    string
	Payload []byte
}

const (
	UserCreate = "user_create"
	UserDelete = "user_delete"
	UserUpdateProfile = "user_update_profile"
)

type UserPayload struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"isVerified"`
	CreatedAt  time.Time `json:"createdAt"`
}
