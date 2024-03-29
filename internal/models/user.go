package models

import (
	"github.com/google/uuid"
	"time"
)

type RegisterUser struct {
	Username     string `db:"username" json:"username"`
	PasswordHash string `db:"password_hash" json:"-"`
	Email        string `db:"email" json:"email"`
}

type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Email        string    `db:"email" json:"email"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
