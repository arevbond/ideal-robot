package models

import "time"

type History struct {
	ID        int       `db:"id"`
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
	Type      int       `db:"type"`
}

type CreateHistory struct {
	Text      string    `db:"text"`
	CreatedAt time.Time `db:"created_at"`
	Type      int       `db:"type"`
}

type Reminder struct {
	ID       int    `db:"id"`
	Text     string `db:"text"`
	IsDone   bool   `db:"is_done"`
	Priority int    `db:"priority"`
}

type CreateReminder struct {
	Text     string `db:"text"`
	IsDone   bool   `db:"is_done"`
	Priority int    `db:"priority"`
}
