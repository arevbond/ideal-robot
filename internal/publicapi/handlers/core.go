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

type roomHandler struct {
	service *services.Service
	log     *slog.Logger
}

func RoomRoutes(log *slog.Logger, db storage.Storage, cfg config.MQTTConfig) chi.Router {
	r := chi.NewRouter()
	roomHandler := &roomHandler{services.New(log, db, cfg), log}
	r.Get("/", roomHandler.Rooms)
	r.Post("/", roomHandler.CreateRoom)
	//r.Route("/{id}", func(r chi.Router) {
	//	r.Use(roomHandler.RoomCtx)
	//	r.Get("/", roomHandler.GetRoom)
	//	r.Post("/device", roomHandler.CreateDevice)
	//})
	return r
}

//func (h *roomHandler) CreateDevice(w http.ResponseWriter, r *http.Request) {
//	r.ParseForm()
//	name := r.Form.Get("name")
//	category := r.Form.Get("type")
//	writeTopic := r.Form.Get("write_topic")
//	readTopic := r.Form.Get("read_topic")
//	strRoomID := r.Form.Get("room_id")
//	roomID, err := strconv.Atoi(strRoomID)
//	if err != nil {
//		h.log.Error("can't convert room_id to int", slog.String("room_id", strRoomID))
//		return
//	}
//
//	err = h.service.CreateDevice(roomID, name, writeTopic, readTopic, getCategory(category))
//	if err != nil {
//		h.log.Error("can't create deivce", slog.Any("error", err))
//		return
//	}
//}

func getCategory(category string) int {
	switch category {
	case "temperature":
		return 1
	case "humidity":
		return 2
	case "motion":
		return 3
	}
	return 0
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
	h.ViewRooms(w, r, viewRoomsProp{rooms, []*models.DeviceWithData{}})
}

func (h *roomHandler) Rooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.service.Rooms()
	if err != nil {
		h.log.Error("failed to increment", slog.Any("error", err))
		http.Error(w, "failed to get rooms", http.StatusInternalServerError)
		return
	}
	devices, err := h.service.GetDevices()
	h.ViewRooms(w, r, viewRoomsProp{rooms: rooms, devices: devices})
}

type viewRoomsProp struct {
	rooms   []*models.Room
	devices []*models.DeviceWithData
}

func (h *roomHandler) ViewRooms(w http.ResponseWriter, r *http.Request, props viewRoomsProp) {
	components.Rooms(props.rooms, props.devices).Render(r.Context(), w)
}

type viewRoomProp struct {
	room *models.Room
}

//func (h *roomHandler) ViewRoom(w http.ResponseWriter, r *http.Request, props viewRoomProp) {
//	components.Room(props.room).Render(r.Context(), w)
//}
