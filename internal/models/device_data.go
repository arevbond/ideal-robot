package models

import "time"

type CreateDeviceData struct {
	DeviceID   int         `json:"device_id" db:"device_id"`
	Value      interface{} `json:"value"`
	Unit       string      `json:"unit" db:"unit"`
	ReceivedAt time.Time   `json:"received_at" db:"received_at"`
}

type DeviceData struct {
	ID         int       `json:"id" db:"id"`
	DeviceID   int       `json:"device_id" db:"device_id"`
	Value      string    `json:"val"`
	Unit       string    `json:"unit" db:"unit"`
	ReceivedAt time.Time `json:"received_at" db:"received_at"`
}

type TemperatureData struct {
	Value float64 `json:"value"`
}

type HumidityData struct {
	Value int `json:"value"`
}
