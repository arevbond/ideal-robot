package models

type CreateDevice struct {
	RoomID     int    `json:"room_id" db:"room_id"`
	Name       string `json:"name" db:"name"`
	Type       int    `json:"type" db:"type"`
	Status     bool   `json:"status" db:"status"`
	WriteTopic string `db:"write_topic"`
	ReadTopic  string `db:"read_topic"`
}

type Device struct {
	ID         int    `json:"id" db:"id"`
	RoomID     int    `json:"room_id" db:"room_id"`
	Name       string `json:"name" db:"name"`
	Type       int    `json:"type" db:"type"`
	Status     bool   `json:"status" db:"status"`
	WriteTopic string `db:"write_topic"`
	ReadTopic  string `db:"read_topic"`
}
