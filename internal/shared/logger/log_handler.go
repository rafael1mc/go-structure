package logger

import (
	"context"
	"gomodel/internal/shared/env"
	"log/slog"
	"os"
)

type LogHandler struct {
	env    *env.Env
	logger *slog.Logger
}

func NewLogHandler(
	env *env.Env,
) *LogHandler {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	return &LogHandler{
		env:    env,
		logger: logger,
	}
}

func (h LogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if level == slog.LevelDebug && !h.env.Environment.IsDebug {
		return false
	}

	return h.logger.Enabled(ctx, level)
}

func (h LogHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.logger.Handler().Handle(ctx, record)
}

func (h LogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.logger.Handler().WithAttrs(attrs)
	return h
}

func (h LogHandler) WithGroup(name string) slog.Handler {
	h.logger.Handler().WithGroup(name)
	return h
}
