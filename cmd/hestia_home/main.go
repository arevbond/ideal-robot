package main

import (
	"HestiaHome/internal/config"
	"log/slog"
	"os"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.New()
	lg := setupLogger(cfg.Env)
	lg = lg.With(slog.String("env", cfg.Env))
	lg.Info("Start server", slog.String("address", cfg.Address))
	lg.Debug("Debug mode enable")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}
