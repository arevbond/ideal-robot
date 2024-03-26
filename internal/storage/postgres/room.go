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

	hubs := []*models.Room{}
	err := s.db.SelectContext(ctx, &hubs, q)
	if err != nil {
		return nil, e.Wrap("cant get all hubs from storage", err)
	}
	return hubs, nil
}

func (s *Storage) CreateRoom(ctx context.Context, hub *models.CreateRoom) (int, error) {
	q1 := `INSERT INTO rooms (name, description) VALUES ($1, $2) RETURNING id`

	var err error
	var rows *sql.Rows

	rows, err = s.db.QueryContext(ctx, q1, hub.Name, hub.Description)

	if err != nil {
		return -1, e.Wrap("can't create hub in storage", err)
	}

	var id int
	if rows.Next() {
		rows.Scan(&id)
	}
	s.log.Debug("create hub", slog.Int("id", id))
	return id, nil
}

func (s *Storage) GetRoomByID(ctx context.Context, id int) (*models.Room, error) {
	q := `SELECT * FROM rooms WHERE id = $1`

	var hub models.Room
	err := s.db.GetContext(ctx, &hub, q, id)
	if err != nil {
		return nil, e.Wrap("can't get hub by id in storage", err)
	}
	return &hub, nil
}

//func (s *Storage) GetRoomsByUserID(ctx context.Context, id uuid.UUID) ([]*models.Room, error) {
//	q := `SELECT * FROM rooms WHERE user_id = $1`
//
//	hubs := []*models.Room{}
//	err := s.db.SelectContext(ctx, &hubs, q, id)
//	if err != nil {
//		return nil, e.Wrap("can't get hubs by user_id in storage", err)
//	}
//	return hubs, nil
//}

func (s *Storage) UpdateRoom(ctx context.Context, hub *models.Room) error {
	q := `UPDATE rooms SET name = $1, description = $2 WHERE id = $3`

	_, err := s.db.ExecContext(ctx, q, hub.Name, hub.Description, hub.ID)
	if err != nil {
		return e.Wrap("can't update hub in storage", err)
	}
	s.log.Debug("update hub", slog.Int("id", hub.ID))
	return nil
}

func (s *Storage) DeleteRoom(ctx context.Context, id int) error {
	q := `DELETE FROM rooms WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return e.Wrap("can't delete hub from storage", err)
	}
	s.log.Debug("delete hub", slog.Int("id", id))
	return nil
}
