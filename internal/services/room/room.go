package room

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/storage"
	"HestiaHome/internal/utils/e"
	"context"
	"log/slog"
)

type Service struct {
	log *slog.Logger
	db  storage.Storage
}

func New(log *slog.Logger, db storage.Storage) *Service {
	return &Service{log: log, db: db}
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

func (s *Service) GetRoomByID(id int) (*models.Room, error) {
	room, err := s.db.GetRoomByID(context.Background(), id)
	if err != nil {
		return nil, e.Wrap("can't get room by id", err)
	}
	return room, err
}
