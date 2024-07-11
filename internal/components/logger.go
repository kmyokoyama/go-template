package components

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	env := os.Getenv("SERVICE_ENV")

	var programLevel = new(slog.LevelVar)
	
	if env == "staging" {
		programLevel.Set(slog.LevelDebug)
	} else {
		programLevel.Set(slog.LevelInfo)
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel}))
}