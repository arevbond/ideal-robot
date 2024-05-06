package handlers

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/publicapi/components"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"strconv"
)

func (h *handler) roomsInDashboard(w http.ResponseWriter, r *http.Request) {
	if strID := chi.URLParam(r, "roomID"); strID != "" {
		id, err := strconv.Atoi(strID)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		devices, err := h.service.GetDevicesByRoomID(id)
		if err != nil {
			h.log.Error("can't get devices", "error", err)
			http.Error(w, "can't get devices", http.StatusInternalServerError)
			return
		}

		h.viewDevicesInDashboard(w, r, devices)
		return
	}
	http.Error(w, "empty ID", http.StatusBadRequest)
}

func (h *handler) dashboard(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.AllRooms()
	if err != nil {
		h.log.Error("failed to get rooms", slog.Any("error", err))
		http.Error(w, "failed to get rooms", http.StatusInternalServerError)
		return
	}
	devices, err := h.service.GetAllDevicesWithData()
	if err != nil {
		h.log.Error("failed to get devices", slog.Any("error", err))
		http.Error(w, "failed to get devices", http.StatusInternalServerError)
		return
	}

	histories, err := h.service.GetHistories(5)
	if err != nil {
		h.log.Error("failed to get histories", slog.Any("error", err))
		http.Error(w, "failed to get histories", http.StatusInternalServerError)
		return
	}

	reminders, err := h.service.AllReminders(5)
	if err != nil {
		h.log.Error("failed to get reminders", slog.Any("error", err))
		http.Error(w, "failed to get reminders", http.StatusInternalServerError)
		return
	}

	h.viewDashboardPage(w, r, viewDashboardProp{rooms: rooms, devices: devices,
		histories: histories, reminders: reminders})
}

func (h *handler) powerDevice(w http.ResponseWriter, r *http.Request) {
	if strID := chi.URLParam(r, "id"); strID != "" {
		id, err := strconv.Atoi(strID)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}

		device, err := h.service.PowerDevice(id)
		if err != nil {
			http.Error(w, "can't power device", http.StatusBadRequest)
			return
		}

		h.viewOffButton(w, r, device)
	}
}

func (h *handler) updateReminder(w http.ResponseWriter, r *http.Request) {
	if strID := chi.URLParam(r, "id"); strID != "" {
		id, err := strconv.Atoi(strID)
		if err != nil {
			http.Error(w, "invalid ID", http.StatusBadRequest)
			return
		}

		isDoneStr := r.FormValue("isDone")
		isDone, err := strconv.ParseBool(isDoneStr)
		if err != nil {
			http.Error(w, "invalid bool value", http.StatusBadRequest)
			return
		}

		err = h.service.UpdateReminder(id, isDone)
		if err != nil {
			http.Error(w, "can't update value", http.StatusBadRequest)
			return
		}

		reminders, err := h.service.AllReminders(5)
		if err != nil {
			h.log.Error("failed to get reminders", slog.Any("error", err))
			http.Error(w, "failed to get reminders", http.StatusInternalServerError)
			return
		}

		h.viewReminders(w, r, reminders)
	}
}

func (h *handler) allHistory(w http.ResponseWriter, r *http.Request) {
	histories, err := h.service.GetHistories(5)
	if err != nil {
		h.log.Error("failed to get histories", slog.Any("error", err))
		http.Error(w, "failed to get histories", http.StatusInternalServerError)
		return
	}

	h.viewHistory(w, r, histories)
}

type viewDashboardProp struct {
	rooms     []*models.Room
	devices   []*models.DeviceWithData
	histories []*models.History
	reminders []*models.Reminder
}

func (h *handler) viewDashboardPage(w http.ResponseWriter, r *http.Request, props viewDashboardProp) {
	components.Dashboard(props.rooms, props.devices, props.histories, props.reminders).Render(r.Context(), w)
}

func (h *handler) viewDevicesInDashboard(w http.ResponseWriter, r *http.Request, devices []*models.DeviceWithData) {
	components.DashboardDevices(devices).Render(r.Context(), w)
}

func (h *handler) viewHistory(w http.ResponseWriter, r *http.Request, histories []*models.History) {
	components.History(histories).Render(r.Context(), w)
}

func (h *handler) viewReminders(w http.ResponseWriter, r *http.Request, reminders []*models.Reminder) {
	components.DashboardReminders(reminders).Render(r.Context(), w)
}

func (h *handler) viewOffButton(w http.ResponseWriter, r *http.Request, device *models.DeviceWithData) {
	components.OffButton(device).Render(r.Context(), w)
}
