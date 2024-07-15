package controllers

import (
	"github.com/google/uuid"
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/models"
)

func ProcessWork(c *components.Components, id uuid.UUID, description string) (models.Work, error) {
	c.Logger.Info("received work", "id", id, "description", description)
	return models.Work{}, nil
}