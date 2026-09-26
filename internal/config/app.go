package config

import (
	"log/slog"
	"time"
)

type AppCfg struct {
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

func loadAppCfg(r *envReader) *AppCfg {
	return &AppCfg{
		LogLevel:        r.logLevel("LOG_LEVEL", slog.LevelInfo),
		ShutdownTimeout: r.duration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}
