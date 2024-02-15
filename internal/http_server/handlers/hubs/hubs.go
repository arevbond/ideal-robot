package hubs

import (
	"HestiaHome/internal/database"
	"HestiaHome/internal/lib/api/response"
	"HestiaHome/internal/models"
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
	db  database.Storage
}

func HubRoutes(log *slog.Logger, db database.Storage) chi.Router {
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
		var hub *models.DBHub

		if idStr := chi.URLParam(r, "id"); idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				h.log.Error("can't convert id to int", err)
				render.Render(w, r, response.ErrInvalidParams)
				return
			}
			hub, err = h.db.GetHubByID(context.Background(), id)
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
	hub := data.Hub
	
	id, err := h.db.CreateHub(context.Background(), hub) // FIXME: bug with null owner_id
	if err != nil {
		h.log.Error("can't create hub", err)
		render.Render(w, r, response.ErrStorageMistake)
		return
	}
	render.Status(r, http.StatusCreated)
	render.Render(w, r, h.NewHubResponse(models.NewDBHub(id, hub)))
}

func (h *HubHandler) GetHub(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*models.DBHub)

	if err := render.Render(w, r, h.NewHubResponse(hub)); err != nil {
		render.Render(w, r, response.ErrRender(err))
		return
	}
}
func (h *HubHandler) UpdateHub(w http.ResponseWriter, r *http.Request) {
	oldHub := r.Context().Value("hub").(*models.DBHub)

	data := &HubRequest{Hub: &models.Hub{OwnerID: oldHub.OwnerID,
		Name: oldHub.Name, Description: oldHub.Description}}
	if err := render.Bind(r, data); err != nil {
		h.log.Error("invalid params for request", err)
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	hub := data.Hub
	newHub := models.NewDBHub(oldHub.ID, hub)
	err := h.db.UpdateHub(context.Background(), newHub)
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, h.NewHubResponse(newHub))
}

func (h *HubHandler) DeleteHub(w http.ResponseWriter, r *http.Request) {
	hub := r.Context().Value("hub").(*models.DBHub)

	err := h.db.DeleteHub(context.Background(), hub.ID)
	if err != nil {
		render.Render(w, r, response.ErrInvalidRequest(err))
		return
	}
	render.Render(w, r, h.NewHubResponse(hub))
}

type HubRequest struct {
	*models.Hub
}

func (h *HubRequest) Bind(r *http.Request) error {
	// h.Hub is nil if no Hub fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if h.Hub == nil {
		return errors.New("missing required Hub fields.")
	}

	// just a post-process after a decode
	h.Hub.Name = strings.ToLower(h.Hub.Name) // as an example, we down-case
	return nil
}

type HubResponse struct {
	Hub  *models.DBHub  `json:"hub"`
	User *models.DBUser `json:"user,omitempty"`
}

func (h *HubHandler) NewHubResponse(hub *models.DBHub) *HubResponse {
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
