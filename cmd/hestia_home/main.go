package main

import (
	"HestiaHome/internal/config"
	"HestiaHome/internal/database/postgres"
	"HestiaHome/internal/models"
	"context"
	"log/slog"
	"os"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.New()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("Start server", slog.String("address", cfg.Address))
	log.Debug("Debug mode enable")

	db, err := postgres.New(log, cfg)
	if err != nil {
		log.Error("Can't init database", err)
		os.Exit(1)
	}
	log.Info("Success connect to database")

	user := &models.DBUser{
		Username:     "Nikita",
		PasswordHash: "123456",
		Email:        "email@gmail.com",
		CreatedAt:    time.Now(),
	}
	err = db.CreateUser(context.Background(), user)
	if err != nil {
		log.Error("can't create user", err)
	} else {
		log.Info("Success create user")
	}

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
