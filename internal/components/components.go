package components

import (
	log "log/slog"
)

type Components struct {
	Logger *log.Logger
	Db     Database
}

func NewComponents(logger *log.Logger, db Database) *Components {
	return &Components{Logger: logger, Db: db}
}
