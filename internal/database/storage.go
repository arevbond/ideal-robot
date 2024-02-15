package database

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
	CreateHub(ctx context.Context, hub *models.Hub) (int, error)
	GetHubByID(ctx context.Context, id int) (*models.DBHub, error)
	GetHubsByUserID(ctx context.Context, id uuid.UUID) ([]*models.DBHub, error)
	UpdateHub(ctx context.Context, hub *models.DBHub) error
	DeleteHub(ctx context.Context, id int) error

	GetUserByID(ctx context.Context, id uuid.UUID) (*models.DBUser, error)
}
