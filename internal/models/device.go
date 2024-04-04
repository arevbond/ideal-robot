package models

type Device struct {
	ID       int    `json:"id" db:"id"`
	RoomID   *int   `json:"room_id" db:"room_id"`
	Name     string `json:"name" db:"name"`
	Category int    `json:"category" db:"category"`
	Hidden   bool   `json:"hidden" db:"hidden"`
	Status   bool   `json:"status" db:"status"`
}
