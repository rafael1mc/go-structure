package logger

import (
	"log/slog"
)

func NewLogger(
	logHandler *LogHandler,
) *slog.Logger {
	return slog.New(logHandler)
}
