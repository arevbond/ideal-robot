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
	//UserRepository
	RoomRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.RegisterUser) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) error
}

type RoomRepository interface {
	GetRooms(ctx context.Context) ([]*models.Room, error)
	CreateRoom(ctx context.Context, room *models.CreateRoom) (int, error)
	GetRoomByID(ctx context.Context, id int) (*models.Room, error)
	UpdateRoom(ctx context.Context, room *models.Room) error
	DeleteRoom(ctx context.Context, id int) error
}
