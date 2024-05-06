package handlers

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/publicapi/components"
	"log/slog"
	"net/http"
)

func (h *handler) Reminders(w http.ResponseWriter, r *http.Request) {
	reminders, err := h.service.AllReminders(100)
	if err != nil {
		h.log.Error("failed to get reminders", slog.Any("error", err))
		http.Error(w, "failed to get reminders", http.StatusInternalServerError)
		return
	}

	h.viewRemindersPage(w, r, reminders)
}

func (h *handler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if text := r.Form.Get("text"); text != "" {
		priority := r.Form.Get("priority")
		allReminders, err := h.service.CreateReminder(text, priority)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		h.viewAllReminders(w, r, allReminders)
		return
	}
	http.Error(w, "empty text", http.StatusBadRequest)
}

func (h *handler) viewRemindersPage(w http.ResponseWriter, r *http.Request, reminders []*models.Reminder) {
	components.Reminders(reminders).Render(r.Context(), w)
}

func (h *handler) viewAllReminders(w http.ResponseWriter, r *http.Request, reminders []*models.Reminder) {
	components.AllReminders(reminders).Render(r.Context(), w)
}
