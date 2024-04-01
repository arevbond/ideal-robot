package postgres

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/utils/e"
	"context"
	"database/sql"
	"log/slog"
)

func (s *Storage) GetRooms(ctx context.Context) ([]*models.Room, error) {
	q := `SELECT * from rooms`

	rooms := []*models.Room{}
	err := s.db.SelectContext(ctx, &rooms, q)
	if err != nil {
		return nil, e.Wrap("cant get all rooms from storage", err)
	}
	return rooms, nil
}

func (s *Storage) CreateRoom(ctx context.Context, room *models.CreateRoom) (int, error) {
	q1 := `INSERT INTO rooms (name, description) VALUES ($1, $2) RETURNING id`

	var err error
	var rows *sql.Rows

	rows, err = s.db.QueryContext(ctx, q1, room.Name, room.Description)

	if err != nil {
		return -1, e.Wrap("can't create room in storage", err)
	}

	var id int
	if rows.Next() {
		rows.Scan(&id)
	}
	s.log.Debug("create room", slog.Int("id", id))
	return id, nil
}

func (s *Storage) GetRoomByID(ctx context.Context, id int) (*models.Room, error) {
	q := `SELECT * FROM rooms WHERE id = $1`

	var room models.Room
	err := s.db.GetContext(ctx, &room, q, id)
	if err != nil {
		return nil, e.Wrap("can't get room by id in storage", err)
	}
	return &room, nil
}

func (s *Storage) UpdateRoom(ctx context.Context, room *models.Room) error {
	q := `UPDATE rooms SET name = $1, description = $2 WHERE id = $3`

	_, err := s.db.ExecContext(ctx, q, room.Name, room.Description, room.ID)
	if err != nil {
		return e.Wrap("can't update room in storage", err)
	}
	s.log.Debug("update room", slog.Int("id", room.ID))
	return nil
}

func (s *Storage) DeleteRoom(ctx context.Context, id int) error {
	q := `DELETE FROM rooms WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return e.Wrap("can't delete room from storage", err)
	}
	s.log.Debug("delete room", slog.Int("id", id))
	return nil
}
