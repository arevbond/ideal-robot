package handlers

import (
	"HestiaHome/internal/config"
	"HestiaHome/internal/services"
	"HestiaHome/internal/storage"
	"github.com/go-chi/chi/v5"
	"log/slog"
)

type handler struct {
	service *services.Service
	log     *slog.Logger
}

func Routes(log *slog.Logger, db storage.Storage, cfg config.MQTTConfig) chi.Router {
	r := chi.NewRouter()
	h := &handler{services.New(log, db, cfg), log}

	r.Get("/", h.dashboard)

	r.Get("/dashboard/{roomID}", h.roomsInDashboard)

	r.Route("/device", func(r chi.Router) {
		r.Post("/power/{id}", h.powerDevice)
	})

	r.Get("/history", h.allHistory)

	r.Route("/reminder", func(r chi.Router) {
		r.Post("/", h.CreateReminder)
		r.Get("/", h.Reminders)
		r.Post("/{id}", h.updateReminder)
	})

	r.Route("/room/{id}", func(r chi.Router) {
		r.Delete("/", h.deleteRoom)
	})
	return r
}
