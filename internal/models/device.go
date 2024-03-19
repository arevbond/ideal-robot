package models

type CreateDevice struct {
	RoomID int    `json:"room_id" db:"room_id"`
	Name   string `json:"name" db:"name"`
	Type   int    `json:"type" db:"type"`
	Status bool   `json:"status" db:"status"`
}

type Device struct {
	ID     int    `json:"id" db:"id"`
	RoomID int    `json:"room_id" db:"room_id"`
	Name   string `json:"name" db:"name"`
	Type   int    `json:"type" db:"type"`
	Status bool   `json:"status" db:"status"`
}
