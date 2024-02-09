package models

type Device struct {
	HubID    int    `json:"hub_id" db:"hub_id"`
	Name     string `json:"name" db:"name"`
	Type     int    `json:"type" db:"type"`
	Location string `json:"location" db:"location"`
	Status   bool   `json:"status" db:"status"`
}

type DBDevice struct {
	ID       int    `json:"id" db:"id"`
	HubID    int    `json:"hub_id" db:"hub_id"`
	Name     string `json:"name" db:"name"`
	Type     int    `json:"type" db:"type"`
	Location string `json:"location" db:"location"`
	Status   bool   `json:"status" db:"status"`
}
