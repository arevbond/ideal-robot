package postgres

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/utils/e"
	"context"
)

func (s *Storage) GetReminders(ctx context.Context, limit int) ([]*models.Reminder, error) {
	q := `SELECT * FROM reminders ORDER BY priority LIMIT $1`

	result := []*models.Reminder{}
	err := s.db.SelectContext(ctx, &result, q, limit)
	if err != nil {
		s.log.Error("can't get all reminders", "error", err)
		return nil, e.Wrap("can't getl all reminders", err)
	}
	return result, nil
}

func (s *Storage) CreateReminder(ctx context.Context, reminder *models.CreateReminder) error {
	q := `INSERT INTO reminders (text, is_done, priority) VALUES ($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, q, reminder.Text, reminder.IsDone, reminder.Priority)
	if err != nil {
		s.log.Error("can't create reminder", "error", err)
		return e.Wrap("can't create reminder", err)
	}

	return nil
}

func (s *Storage) UpdateReminder(ctx context.Context, reminder *models.Reminder) error {
	q := `UPDATE reminders SET text = $1, is_done = $2, priority = $3 WHERE id = $4`

	_, err := s.db.ExecContext(ctx, q, reminder.Text, reminder.IsDone,
		reminder.Priority, reminder.ID)
	if err != nil {
		s.log.Error("can't update reminder", "error", err)
		return e.Wrap("can't update reminder", err)
	}
	return nil
}

func (s *Storage) DeleteReminder(ctx context.Context, id int) error {
	q := `DELETE FROM reminders WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		s.log.Error("can't delete reminders", "error", err)
		return e.Wrap("can't delete reminders", err)
	}

	return nil
}
