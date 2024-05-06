package services

import (
	"HestiaHome/internal/clients/mqtt"
	"HestiaHome/internal/config"
	"HestiaHome/internal/models"
	"HestiaHome/internal/storage"
	"HestiaHome/internal/utils/e"
	"context"
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"strconv"
)

type Service struct {
	log        *slog.Logger
	db         storage.Storage
	mqttClient mqtt2.Client
}

func New(log *slog.Logger, db storage.Storage, cfg config.MQTTConfig) *Service {
	client := mqtt.New(cfg.Address, cfg.Port, cfg.ClientID, cfg.Username, cfg.Password)
	mqtt.Subscribe("topic/test", client)
	s := &Service{log: log, db: db, mqttClient: client}
	go processDataFromMQTT(log, db, s)
	return s
}

func (s *Service) AllRooms() ([]*models.Room, error) {
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

func (s *Service) GetDevicesByRoomID(id int) ([]*models.DeviceWithData, error) {
	devices, err := s.db.GetDevicesWithDataByRoomID(context.Background(), id)
	if err != nil {
		return nil, e.Wrap("service can't get devices by room id", err)
	}
	return devices, nil
}

func (s *Service) GetAllDevicesWithData() ([]*models.DeviceWithData, error) {
	devices, err := s.db.GetDevicesWithData(context.Background())
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

func (s *Service) PowerDevice(id int) (*models.DeviceWithData, error) {
	device, err := s.db.GetDeviceByID(context.Background(), id)
	if err != nil {
		return nil, e.Wrap("can't get device", err)
	}
	device.Status = !device.Status

	mqtt.Publish(fmt.Sprintf("state/%d", id), s.mqttClient, strconv.FormatBool(device.Status))
	s.log.Debug("send messega in topic", slog.Int("id", id))

	err = s.db.UpdateDevice(context.Background(), device)
	if err != nil {
		return nil, e.Wrap("can't update device", err)
	}

	result, err := s.db.GetDevicesWithDataByID(context.Background(), id)
	if err != nil {
		return nil, e.Wrap("can't get device with data from db", err)
	}

	return result, nil
}
