package wire

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name string `json:"name"`
}

type UserResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
