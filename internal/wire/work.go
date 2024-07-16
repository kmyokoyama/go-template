package wire

import "github.com/google/uuid"

type WorkRequest struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

type WorkResponse struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}
