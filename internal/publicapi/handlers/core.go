package handlers

import (
	"HestiaHome/internal/config"
	"HestiaHome/internal/models"
	"HestiaHome/internal/publicapi/components"
	"HestiaHome/internal/services"
	"HestiaHome/internal/storage"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type handler struct {
	service *services.Service
	log     *slog.Logger
}

func Routes(log *slog.Logger, db storage.Storage, cfg config.MQTTConfig) chi.Router {
	r := chi.NewRouter()
	roomHandler := &handler{services.New(log, db, cfg), log}
	r.Get("/", roomHandler.Dashboard)
	return r
}

func (h *handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.Rooms()
	if err != nil {
		h.log.Error("failed to increment", slog.Any("error", err))
		http.Error(w, "failed to get rooms", http.StatusInternalServerError)
		return
	}
	devices, err := h.service.GetDevices()
	if err != nil {
		h.log.Error("failed to increment", slog.Any("error", err))
		http.Error(w, "failed to get rooms", http.StatusInternalServerError)
		return
	}
	h.viewDashboard(w, r, viewDashboardProp{rooms: rooms, devices: devices})
}

type viewDashboardProp struct {
	rooms   []*models.Room
	devices []*models.DeviceWithData
}

func (h *handler) viewDashboard(w http.ResponseWriter, r *http.Request, props viewDashboardProp) {
	components.Dashboard(props.rooms, props.devices).Render(r.Context(), w)
}
