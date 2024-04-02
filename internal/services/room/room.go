package room

import (
	"HestiaHome/internal/clients/mqtt"
	"HestiaHome/internal/config"
	"HestiaHome/internal/models"
	"HestiaHome/internal/storage"
	"HestiaHome/internal/utils/e"
	"context"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
)

type Service struct {
	log        *slog.Logger
	db         storage.Storage
	mqttClient mqtt2.Client
}

func New(log *slog.Logger, db storage.Storage, cfg config.MQTTConfig) *Service {
	client := mqtt.New(cfg.Address, cfg.Port, cfg.ClientID, cfg.Username, cfg.Password)
	return &Service{log: log, db: db, mqttClient: client}
}

func (s *Service) Rooms() ([]*models.Room, error) {
	rooms, err := s.db.GetRooms(context.Background())
	if err != nil {
		return nil, e.Wrap("can't get rooms", err)
	}
	return rooms, nil
}

func (s *Service) CreateRoom(name string) error {
	createRoom := &models.CreateRoom{Name: name}
	_, err := s.db.CreateRoom(context.Background(), createRoom)
	if err != nil {
		return e.Wrap("can't create room", err)
	}
	return nil
}

func (s *Service) GetRoom(id int) (*models.Room, error) {
	room, err := s.db.GetRoomByID(context.Background(), id)
	if err != nil {
		return nil, e.Wrap("can't get room by id", err)
	}
	return room, err
}

func (s *Service) CreateDevice(roomID int, name, writeTopic, readTopic string) error {
	err := s.db.CreateDevice(context.Background(), &models.CreateDevice{
		RoomID:     roomID,
		Name:       name,
		WriteTopic: writeTopic,
		ReadTopic:  readTopic,
	})
	if err != nil {
		return e.Wrap("can't create device", err)
	}
	return nil
}

func (s *Service) GetDevices(id int) ([]*models.Device, error) {
	devices, err := s.db.GetDevicesByRoomID(context.Background(), id, 0, 10)
	if err != nil {
		return nil, e.Wrap("service can't get devices by room id", err)
	}
	return devices, nil
}

func (s *Service) DeleteRoom(id int) error {
	err := s.db.DeleteRoom(context.Background(), id)
	if err != nil {
		return e.Wrap("can't delete room by id", err)
	}
	return nil
}
