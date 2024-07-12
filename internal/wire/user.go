package wire

import "github.com/google/uuid"

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserResponse struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role string    `json:"role"`
}
