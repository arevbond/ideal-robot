package services

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/utils/e"
	"context"
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

func (s *Service) CreateReminders(text string) error {
	createReminder := &models.CreateReminder{Text: text, IsDone: false, Priority: 0}
	err := s.db.CreateReminder(context.Background(), createReminder)
	if err != nil {
		return e.Wrap("can't create reminder", err)
	}
	return nil
}
