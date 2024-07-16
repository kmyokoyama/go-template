package components

import (
	log "log/slog"
)

type Components struct {
	Logger *log.Logger
	Config Config
	Db     Database
}

func NewComponents(logger *log.Logger, config Config, db Database) *Components {
	return &Components{Logger: logger, Config: config, Db: db}
}
