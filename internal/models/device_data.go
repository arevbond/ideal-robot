package models

import "time"

type DeviceData struct {
	DeviceID   int       `json:"device_id" db:"device_id"`
	Value      string    `json:"val"`
	Unit       string    `json:"unit" db:"unit"`
	ReceivedAt time.Time `json:"received_at" db:"received_at"`
}

type DBDeviceData struct {
	ID         int       `json:"id" db:"id"`
	DeviceID   int       `json:"device_id" db:"device_id"`
	Value      string    `json:"val"`
	Unit       string    `json:"unit" db:"unit"`
	ReceivedAt time.Time `json:"received_at" db:"received_at"`
}
