package controllers

import (
	"github.com/google/uuid"
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/models"
)

func CreateUser(c *components.Components, user models.User) (models.User, error) {
	id := uuid.New()
	user.Id = id

	err := c.Db.CreateUser(user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func FindUser(c *components.Components, id uuid.UUID) (models.User, error) {
	c.Logger.Info("user.go:", "id", id)
	user, err := c.Db.FindUser(id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}