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
	"strconv"
)

type handler struct {
	service *services.Service
	log     *slog.Logger
}

func Routes(log *slog.Logger, db storage.Storage, cfg config.MQTTConfig) chi.Router {
	r := chi.NewRouter()
	h := &handler{services.New(log, db, cfg), log}
	r.Get("/{roomID}", h.Dashboard)
	r.Get("/", h.Dashboard)

	r.Route("/room/{id}", func(r chi.Router) {
		r.Delete("/", h.DeleteRoom)
	})
	return r
}

func (h *handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	if strID := chi.URLParam(r, "roomID"); strID != "" {
		id, err := strconv.Atoi(strID)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		room, err := h.service.GetRoom(id)
		if err != nil {
			http.Error(w, "can't get room", http.StatusBadRequest)
			return
		}

		devices, err := h.service.GetDevicesByRoomID(id)
		if err != nil {
			h.log.Error("can't get devices", "error", err)
			http.Error(w, "can't get devices", http.StatusBadRequest)
			return
		}

		h.viewDevicesInDashboard(w, r, viewDashboardProp{rooms: []*models.Room{room}, devices: devices})
		return
	}

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

func (h *handler) viewDevicesInDashboard(w http.ResponseWriter, r *http.Request, props viewDashboardProp) {
	components.DashboardDevices(props.devices).Render(r.Context(), w)
}
