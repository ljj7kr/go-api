package logger

import (
	"log/slog"
	"os"

	"go-api/internal/config"
)

func New(cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.LogLevel(),
	}

	// 로컬 개발에서는 사람이 읽기 쉬운 text format 사용
	if cfg.AppEnv == "local" || cfg.AppEnv == "dev" {
		return slog.New(slog.NewTextHandler(os.Stdout, opts))
	}

	// 그 외 환경에서는 수집기 친화적인 JSON format 사용
	return slog.New(slog.NewJSONHandler(os.Stdout, opts))
}
