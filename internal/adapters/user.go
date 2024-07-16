package adapters

import (
	"github.com/kmyokoyama/go-template/internal/models"
	"github.com/kmyokoyama/go-template/internal/wire"
)

func ToUserInternal(w wire.SignupRequest) (models.User, error) {
	role, err := models.RoleFromString(w.Role)
	if err != nil {
		return models.User{}, err
	}

	return models.User{Username: w.Username, Role: role}, nil
}

func ToUserResponse(user models.User) wire.UserResponse {
	return wire.UserResponse{Id: user.Id, Username: user.Username, Role: user.Role.String()}
}
