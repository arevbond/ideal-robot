package services

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/utils/e"
	"context"
	"log/slog"
)

func (s *Service) GetReminders(limit int) ([]*models.Reminder, error) {
	if limit < 0 {
		limit = 5
	}
	reminders, err := s.db.GetReminders(context.Background(), limit)
	if err != nil {
		return nil, e.Wrap("can't get histories", err)
	}
	return reminders, nil
}

func (s *Service) GetReminder(id int) (*models.Reminder, error) {
	reminder, err := s.db.GetReminder(context.Background(), id)
	if err != nil {
		return nil, e.Wrap("can't process service get reminder", err)
	}
	return reminder, nil
}

func (s *Service) CreateReminders(text string) error {
	createReminder := &models.CreateReminder{Text: text, IsDone: false, Priority: 0}
	err := s.db.CreateReminder(context.Background(), createReminder)
	if err != nil {
		return e.Wrap("can't create reminder", err)
	}
	return nil
}

func (s *Service) UpdateReminder(id int, isDone bool) error {
	reminder, err := s.db.GetReminder(context.Background(), id)
	if err != nil {
		s.log.Error("can't get reminder from db", slog.Any("error", err))
		return e.Wrap("can't get reminder from db", err)
	}

	reminder.IsDone = isDone

	err = s.db.UpdateReminder(context.Background(), reminder)
	if err != nil {
		s.log.Error("can't update reminder", slog.Any("error", err))
		return e.Wrap("can't update reminder", err)
	}
	return nil
}
