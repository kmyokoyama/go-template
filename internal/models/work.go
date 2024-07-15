package models

import "github.com/google/uuid"

type Work struct {
	Id uuid.UUID
	Description string
	Status string // TODO: Make it safe-enum later.
}