package models

import "time"

type Device struct {
	ID       int    `json:"id" db:"id"`
	RoomID   *int   `json:"room_id" db:"room_id"`
	Name     string `json:"name" db:"name"`
	Category int    `json:"category" db:"category"`
	Hidden   bool   `json:"hidden" db:"hidden"`
	Status   bool   `json:"status" db:"status"`
}

type DeviceWithData struct {
	ID         int       `db:"device_id"`
	Name       string    `db:"device_name"`
	Value      Value     `db:"value"`
	Unit       string    `db:"unit"`
	ReceivedAt time.Time `db:"received_at"`
}

type Value struct {
	Value float64 `db:"value" json:"value"`
}
