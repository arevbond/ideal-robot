package storage

import (
	"HestiaHome/internal/models"
	"context"
	"errors"
	"github.com/google/uuid"
)

var (
	ErrUserNotExist = errors.New("user not exist")
)

type Storage interface {
	CreateRoom(ctx context.Context, hub *models.CreateRoom) (int, error)
	GetRoomByID(ctx context.Context, id int) (*models.Room, error)
	GetRoomsByUserID(ctx context.Context, id uuid.UUID) ([]*models.Room, error)
	UpdateRoom(ctx context.Context, hub *models.Room) error
	DeleteRoom(ctx context.Context, id int) error

	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}
