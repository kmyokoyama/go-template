package components

import (
	"log/slog"
	"os"
)

func NewLogger(config Config) *slog.Logger {
	env := config.Get("SERVICE_ENV")

	var programLevel = new(slog.LevelVar)
	
	if env == "staging" {
		programLevel.Set(slog.LevelDebug)
	} else {
		programLevel.Set(slog.LevelInfo)
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
}