package logger

import (
	"bank_app/internal/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"os"
	"path/filepath"
)

type appLogger struct {
	*slog.Logger
}

func NewLogger(cfg config.Logger) (Logger, error) {
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	if cfg.FilePath != "" {
		logDir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename: cfg.FilePath,
		MaxSize:  cfg.MaxSize,
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler
	handler = slog.NewTextHandler(lumberjackLogger, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &appLogger{Logger: logger}, nil
}
