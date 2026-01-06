package api

import (
	"log/slog"

	ilog "cdpnetool/internal/log"
)

func SetLogger(l *slog.Logger) { ilog.Set(l) }
