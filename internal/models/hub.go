package models

import "github.com/google/uuid"

type Hub struct {
	OwnerID     uuid.UUID `json:"owner_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
}

type DBHub struct {
	ID          int       `json:"-" db:"id"`
	OwnerID     uuid.UUID `json:"owner_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
}
