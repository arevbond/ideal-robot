package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *handler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	if strID := chi.URLParam(r, "id"); strID != "" {
		id, err := strconv.Atoi(strID)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		err = h.service.DeleteRoom(id)
		if err != nil {
			http.Error(w, "Failed to delete room", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	http.Error(w, "Missing ID", http.StatusBadRequest)
}
