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
	DeviceRepository
}

type RoomRepository interface {
	GetRooms(ctx context.Context) ([]*models.Room, error)
	CreateRoom(ctx context.Context, room *models.CreateRoom) (int, error)
	GetRoomByID(ctx context.Context, id int) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id int) error
}

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDeviceByID(ctx context.Context, id int) (*models.Device, error)
	GetDevicesByRoomID(ctx context.Context, roomID int, offset, limit int) ([]*models.Device, error)
}
