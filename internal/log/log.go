package log

import (
	"log/slog"
	"os"
)

var defaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func Set(l *slog.Logger) { defaultLogger = l }

func L() *slog.Logger { return defaultLogger }
