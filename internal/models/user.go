package models

import (
	"github.com/google/uuid"
	"time"
)

type DBUser struct {
	ID           uuid.UUID `db:"id" json:"-"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Email        string    `db:"email" json:"email"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
