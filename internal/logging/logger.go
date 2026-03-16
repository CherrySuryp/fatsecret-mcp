package logging

import (
	"log/slog"
	"os"
)

func NewLogger(Loglevel slog.Level) *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: Loglevel,
			},
		),
	)
}

