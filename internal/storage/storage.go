package storage

import (
	"HestiaHome/internal/models"
	"context"
	"errors"
)

var (
	ErrUserNotExist = errors.New("user not exist")
)

type Storage interface {
	RoomRepository
}

type RoomRepository interface {
	GetRooms(ctx context.Context) ([]*models.Room, error)
	CreateRoom(ctx context.Context, room *models.CreateRoom) (int, error)
	GetRoomByID(ctx context.Context, id int) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id int) error
}

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.CreateDevice) error
	GetDeviceByID(ctx context.Context, id int) (*models.Device, error)
}
