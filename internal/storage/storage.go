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
	HistoryRepository
	ReminderRepository
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
	GetDevicesWithDataByID(ctx context.Context, id int) (*models.DeviceWithData, error)
	GetDevicesWithDataByRoomID(ctx context.Context, id int) ([]*models.DeviceWithData, error)
	GetDevices(ctx context.Context) ([]*models.Device, error)
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDeviceByID(ctx context.Context, id int) (*models.Device, error)
	UpdateDevice(ctx context.Context, device *models.Device) error
	GetDevicesByRoomID(ctx context.Context, roomID int, offset, limit int) ([]*models.Device, error)
}

type DeviceDataRepository interface {
	CreateDeviceData(ctx context.Context, deviceData *models.CreateDeviceData) error
	GetDeviceDataByID(ctx context.Context, id int) (*models.DeviceData, error)
}

type HistoryRepository interface {
	GetHistories(ctx context.Context, limit int) ([]*models.History, error)
	CreateHistory(ctx context.Context, history *models.CreateHistory) error
	UpdateHistory(ctx context.Context, history *models.History) error
	DeleteHistory(ctx context.Context, id int) error
}

type ReminderRepository interface {
	GetReminders(ctx context.Context, limit int) ([]*models.Reminder, error)
	CreateReminder(ctx context.Context, reminder *models.CreateReminder) error
	UpdateReminder(ctx context.Context, reminder *models.Reminder) error
	GetReminder(ctx context.Context, id int) (*models.Reminder, error)
}
