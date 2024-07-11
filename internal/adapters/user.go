package adapters

import (
	"github.com/kmyokoyama/go-template/internal/models"
	"github.com/kmyokoyama/go-template/internal/wire"
)

func ToUserResponse(user models.User) wire.UserResponse {
	return wire.UserResponse{Id: user.Id, Name: user.Name}
}
