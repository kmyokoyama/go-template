package models

import (
	"errors"

	"github.com/google/uuid"
)

type Status struct {
	slug string
}

func (s Status) String() string {
	return s.slug
}

var (
	StatusUnknown Status = Status{"unknown"}
	StatusPending        = Status{"pending"}
	StatusDone           = Status{"done"}
)

func StatusFromString(s string) (Status, error) {
	switch s {
	case StatusPending.slug:
		return StatusPending, nil
	case StatusDone.slug:
		return StatusDone, nil
	default:
		return StatusUnknown, errors.New("unknown status: " + s)
	}
}

type Work struct {
	Id          uuid.UUID
	Description string
	Status      Status
}
