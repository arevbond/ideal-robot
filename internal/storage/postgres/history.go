package postgres

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/utils/e"
	"context"
)

func (s *Storage) GetHistories(ctx context.Context, limit int) ([]*models.History, error) {
	q := `SELECT * FROM history ORDER BY created_at LIMIT $1`

	result := []*models.History{}
	err := s.db.SelectContext(ctx, &result, q, limit)
	if err != nil {
		s.log.Error("can't get all histories", "error", err)
		return nil, e.Wrap("can't getl all histories", err)
	}
	return result, nil
}

func (s *Storage) CreateHistory(ctx context.Context, history *models.CreateHistory) error {
	q := `INSERT INTO history (text, created_at, type) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, q, history.Text, history.CreatedAt, history.Type)
	if err != nil {
		s.log.Error("can't create history", "error", err)
		return e.Wrap("can't create history", err)
	}

	return nil
}

func (s *Storage) UpdateHistory(ctx context.Context, history *models.History) error {
	q := `UPDATE history SET text = $1, createad_at = $2, type = $3 WHERE id = $4`

	_, err := s.db.ExecContext(ctx, q, history.Text, history.CreatedAt, history.Type, history.ID)
	if err != nil {
		s.log.Error("can't update history", "error", err)
		return e.Wrap("can't update history", err)
	}
	return nil
}

func (s *Storage) DeleteHistory(ctx context.Context, id int) error {
	q := `DELETE FROM history WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		s.log.Error("can't delete history", "error", err)
		return e.Wrap("can't delete history", err)
	}

	return nil
}
