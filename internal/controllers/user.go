package controllers

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/models"
	"golang.org/x/crypto/bcrypt"
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

func Login(c *components.Components, username string, password string) (string, error) {
	user, hashedPassword, err := c.Db.FindUserAndPasswordByUsername(username)
	if err != nil {
		c.Logger.Error("login error", "err", err)
		return "", errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		c.Logger.Error("login error", "err", err)
		return "", errors.New("username or password does not match")
	}

	token, err := newToken(user, "secret")
	if err != nil {
		c.Logger.Error("login failed", "err", err)
		return "", err
	}

	return token, nil
}

func FindUser(c *components.Components, id uuid.UUID) (models.User, error) {
	user, err := c.Db.FindUser(id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// JWT.

type UserClaims struct {
	UserId uuid.UUID `json:"user-id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func newToken(user models.User, secret string) (string, error) {
	// TODO: Add time-related claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &UserClaims{
		UserId: user.Id,
		Role:    user.Role.String(),
	})

	// Sign and get the complete encoded token as a string using the secret
	signedToken, err := token.SignedString([]byte(secret)) // TODO: Handle this error.
	if err != nil {
		return "", err
	}

	return signedToken, err
}

func IsValidToken(accessToken string) bool {
	parsedAccessToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return false
	}

	return parsedAccessToken.Valid
}