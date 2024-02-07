package database

import (
	"HestiaHome/internal/models"
	"context"
	"errors"
)

var (
	ErrUserNotExist = errors.New("user not exist")
)

type Database interface {
	CreateUser(ctx context.Context, user *models.DBUser) error
	UpdateUser(ctx context.Context, user *models.DBUser) error
}
