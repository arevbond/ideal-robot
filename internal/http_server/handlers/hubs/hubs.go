package hubs

import (
	"HestiaHome/internal/lib/api/response"
	"HestiaHome/internal/models"
	"HestiaHome/internal/storage"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type HubHandler struct {
	log *slog.Logger
	db  storage.Storage
}

func HubRoutes(log *slog.Logger, db storage.Storage) chi.Router {
	r := chi.NewRouter()
	hubHandler := &HubHandler{log, db}
	r.Post("/", hubHandler.CreateHub)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(hubHandler.HubCtx)
		r.Get("/", hubHandler.GetHub)
		r.Put("/", hubHandler.UpdateHub)
		r.Delete("/", hubHandler.DeleteHub)
	})
	return r
}

func (h *HubHandler) HubCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var hub *models.Room

		if idStr := chi.URLParam(r, "id"); idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				h.log.Error("can't convert id to int", err)
				render.Render(w, r, response.ErrInvalidParams) //nolint:errcheck
				return
			}
			hub, err = h.db.GetRoomByID(context.Background(), id)
			if err != nil {
				render.Render(w, r, response.ErrNotFound)
				return
			}
		} else {
			render.Render(w, r, response.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "hub", hub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *HubHandler) CreateHub(w http.ResponseWriter, r *http.Request) {
	data := &HubRequest{}
	if err := render.Bind(r, data); err != nil {
		h.log.Error("invalid params for request", err)
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	hub := data.CreateRoom

	id, err := h.db.CreateRoom(context.Background(), hub)
	if err != nil {
		h.log.Error("can't create hub", err)
		render.Render(w, r, response.ErrStorageMistake)
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, h.NewHubResponse(models.NewRoom(id, hub)))
}

func (h *HubHandler) GetHub(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*models.Room)

	if err := render.Render(w, r, h.NewHubResponse(hub)); err != nil {
		render.Render(w, r, response.ErrRender(err))
		return
	}
}

func (h *HubHandler) UpdateHub(w http.ResponseWriter, r *http.Request) {
	oldHub := r.Context().Value("hub").(*models.Room)

	data := &HubRequest{CreateRoom: &models.CreateRoom{OwnerID: oldHub.OwnerID,
		Name: oldHub.Name, Description: oldHub.Description}}
	if err := render.Bind(r, data); err != nil {
		h.log.Error("invalid params for request", err)
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	hub := data.CreateRoom
	newHub := models.NewRoom(oldHub.ID, hub)
	err := h.db.UpdateRoom(context.Background(), newHub)
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, h.NewHubResponse(newHub))
}

func (h *HubHandler) DeleteHub(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*models.Room)

	err := h.db.DeleteRoom(context.Background(), hub.ID)
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, h.NewHubResponse(hub))
}

type HubRequest struct {
	*models.CreateRoom
}

func (h *HubRequest) Bind(r *http.Request) error {
	// h.CreateRoom is nil if no CreateRoom fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if h.CreateRoom == nil || h.CreateRoom.Name == "" {
		return errors.New("missing required CreateRoom fields.")
	}

	// just a post-process after a decode
	h.CreateRoom.Name = strings.ToLower(h.CreateRoom.Name) // as an example, we down-case
	return nil
}

type HubResponse struct {
	Hub  *models.Room `json:"hub"`
	User *models.User `json:"user,omitempty"`
}

func (h *HubHandler) NewHubResponse(hub *models.Room) *HubResponse {
	resp := &HubResponse{Hub: hub}
	user, err := h.db.GetUserByID(context.Background(), hub.OwnerID)
	if err == nil {
		resp.User = user
	}
	return resp
}

func (hr *HubResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
