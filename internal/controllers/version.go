package controllers

import (
	"github.com/kmyokoyama/go-template/internal/components"
	"github.com/kmyokoyama/go-template/internal/models"
)

func GetVersion(c *components.Components) (models.Version, error) {
	version, err := c.Db.FindVersion()
	
	if err != nil {
		return models.Version{}, err
	}

	return version, nil
}