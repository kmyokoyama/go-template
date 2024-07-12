package controllers

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/models"
)

func Signup(c *components.Components, user models.User, password string) (models.User, error) {
	id := uuid.New()
	user.Id = id

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	err := c.Db.CreateUser(user, string(hashedPassword))
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
