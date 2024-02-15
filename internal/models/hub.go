package models

import "github.com/google/uuid"

type Hub struct {
	OwnerID     uuid.UUID `json:"owner_id,omitempty" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
}

type DBHub struct {
	ID          int       `json:"id" db:"id"`
	OwnerID     uuid.UUID `json:"owner_id,omitempty" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
}

func NewDBHub(id int, hub *Hub) *DBHub {
	return &DBHub{
		ID:          id,
		OwnerID:     hub.OwnerID,
		Name:        hub.Name,
		Description: hub.Description,
	}
}
