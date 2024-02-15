package main

import (
	"HestiaHome/internal/config"
	"HestiaHome/internal/database/postgres"
	"HestiaHome/internal/http_server/handlers/hubs"
	mwLoger "HestiaHome/internal/http_server/middleware/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"os"
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

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(mwLoger.New(log))
	router.Use(render.SetContentType(render.ContentTypeJSON))

	router.Mount("/hubs", hubs.HubRoutes(log, db))

	router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("Hello World"))
		if err != nil {
			log.Error("", err)
		}
	})

	err = http.ListenAndServe(cfg.Address, router)
	if err != nil {
		log.Error("can't init http server", err)
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
