package handlers

import (
	"HestiaHome/internal/components"
	"HestiaHome/internal/config"
	"HestiaHome/internal/models"
	"HestiaHome/internal/services/room"
	"HestiaHome/internal/storage"
	"HestiaHome/internal/utils/api/response"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type roomHandler struct {
	service *room.Service
	log     *slog.Logger
}

func RoomRoutes(log *slog.Logger, db storage.Storage, cfg config.MQTTConfig) chi.Router {
	r := chi.NewRouter()
	roomHandler := &roomHandler{room.New(log, db, cfg), log}
	r.Get("/", roomHandler.Rooms)
	r.Post("/", roomHandler.CreateRoom)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(roomHandler.RoomCtx)
		r.Get("/", roomHandler.GetRoom)
	})
	return r
}

func (h *roomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	rm := r.Context().Value("room").(*models.Room)

	h.ViewRoom(w, r, viewRoomProp{room: rm})
}

func (h *roomHandler) RoomCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var room *models.Room

		if idStr := chi.URLParam(r, "id"); idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				h.log.Error("can't convert id to int", err)
				render.Render(w, r, response.ErrInvalidParams) //nolint:errcheck
				return
			}
			room, err = h.service.GetRoom(id)
			if err != nil {
				render.Render(w, r, response.ErrNotFound)
				return
			}
		} else {
			render.Render(w, r, response.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "room", room)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *roomHandler) Rooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.Rooms()
	if err != nil {
		h.log.Error("failed to increment", slog.Any("error", err))
		http.Error(w, "failed to get rooms", http.StatusInternalServerError)
		return
	}
	h.ViewRooms(w, r, viewRoomsProp{rooms: rooms})
}

func (h *roomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form.Get("name")
	err := h.service.CreateRoom(name)
	if err != nil {
		h.log.Error("failed to create room", slog.Any("error", err))
		http.Error(w, "failed to create room", http.StatusInternalServerError)
		return
	}
	rooms, err := h.service.Rooms()
	if err != nil {
		h.log.Error("failed to increment", slog.Any("error", err))
		http.Error(w, "failed to get rooms", http.StatusInternalServerError)
		return
	}
	h.ViewRooms(w, r, viewRoomsProp{rooms})
}

type viewRoomsProp struct {
	rooms []*models.Room
}

func (h *roomHandler) ViewRooms(w http.ResponseWriter, r *http.Request, props viewRoomsProp) {
	components.Rooms(props.rooms).Render(r.Context(), w)
}

type viewRoomProp struct {
	room *models.Room
}

func (h *roomHandler) ViewRoom(w http.ResponseWriter, r *http.Request, props viewRoomProp) {
	components.Room(props.room).Render(r.Context(), w)
}
