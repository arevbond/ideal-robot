package postgres

import (
	"HestiaHome/internal/lib/e"
	"HestiaHome/internal/models"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"log/slog"
)

func (s *Storage) GetHubs(ctx context.Context) ([]*models.DBHub, error) {
	q := `SELECT * from hubs`

	hubs := []*models.DBHub{}
	err := s.db.SelectContext(ctx, &hubs, q)
	if err != nil {
		return nil, e.Wrap("cant get all hubs from storage", err)
	}
	return hubs, nil
}

func (s *Storage) CreateHub(ctx context.Context, hub *models.Hub) (int, error) {
	q1 := `INSERT INTO hubs (user_id, name, description) VALUES ($1, $2, $3) RETURNING id`
	q2 := `INSERT INTO hubs (name, description) VALUES ($1, $2) RETURNING id`

	var err error
	var rows *sql.Rows

	if hub.OwnerID != uuid.Nil {
		rows, err = s.db.QueryContext(ctx, q1, hub.OwnerID, hub.Name, hub.Description)
	} else {
		rows, err = s.db.QueryContext(ctx, q2, hub.Name, hub.Description)
	}
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

func (s *Storage) GetHubByID(ctx context.Context, id int) (*models.DBHub, error) {
	q := `SELECT * FROM hubs WHERE id = $1`

	var hub models.DBHub
	err := s.db.GetContext(ctx, &hub, q, id)
	if err != nil {
		return nil, e.Wrap("can't get hub by id in storage", err)
	}
	return &hub, nil
}

func (s *Storage) GetHubsByUserID(ctx context.Context, id uuid.UUID) ([]*models.DBHub, error) {
	q := `SELECT * FROM hubs WHERE user_id = $1`

	hubs := []*models.DBHub{}
	err := s.db.SelectContext(ctx, &hubs, q, id)
	if err != nil {
		return nil, e.Wrap("can't get hubs by user_id in storage", err)
	}
	return hubs, nil
}

func (s *Storage) UpdateHub(ctx context.Context, hub *models.DBHub) error {
	q := `UPDATE hubs SET user_id = $1, name = $2, description = $3 WHERE id = $4`

	_, err := s.db.ExecContext(ctx, q, hub.OwnerID, hub.Name, hub.Description, hub.ID)
	if err != nil {
		return e.Wrap("can't update hub in storage", err)
	}
	s.log.Debug("update hub", slog.Int("id", hub.ID))
	return nil
}

func (s *Storage) DeleteHub(ctx context.Context, id int) error {
	q := `DELETE FROM hubs WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return e.Wrap("can't delete hub from storage", err)
	}
	s.log.Debug("delete hub", slog.Int("id", id))
	return nil
}
