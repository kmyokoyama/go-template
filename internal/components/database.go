package components

import (
	"github.com/google/uuid"
	"github.com/kmyokoyama/go-template/internal/models"
)

type Database interface {
	FindVersion() (models.Version, error)
	CreateUser(models.User) error
	FindUser(id uuid.UUID) (models.User, error)
}
