package storage

import (
	"HestiaHome/internal/models"
	"context"
	"errors"
)

var (
	ErrUserNotExist   = errors.New("user not exist")
	ErrDeviceNotExist = errors.New("device not exist")
)

type Storage interface {
	RoomRepository
	DeviceRepository
	DeviceDataRepository
}

type RoomRepository interface {
	GetRooms(ctx context.Context) ([]*models.Room, error)
	CreateRoom(ctx context.Context, room *models.CreateRoom) (int, error)
	GetRoomByID(ctx context.Context, id int) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id int) error
}

type DeviceRepository interface {
	GetDevicesWithData(ctx context.Context) ([]*models.DeviceWithData, error)
	GetDevices(ctx context.Context) ([]*models.Device, error)
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDeviceByID(ctx context.Context, id int) (*models.Device, error)
	GetDevicesByRoomID(ctx context.Context, roomID int, offset, limit int) ([]*models.Device, error)
}

type DeviceDataRepository interface {
	CreateDeviceData(ctx context.Context, deviceData *models.CreateDeviceData) error
	GetDeviceDataByID(ctx context.Context, id int) (*models.DeviceData, error)
}
