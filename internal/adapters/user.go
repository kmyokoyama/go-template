package adapters

import (
	"github.com/kmyokoyama/go-template/internal/models"
	"github.com/kmyokoyama/go-template/internal/wire"
)

func ToUserInternal(w wire.CreateUserRequest) (models.User, error) {
	role, err := models.FromString(w.Role)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Username: w.Username, Role: role}, nil
}

func ToUserResponse(user models.User) wire.UserResponse {
	return wire.UserResponse{Id: user.Id, Name: user.Username, Role: user.Role.String()}
}
