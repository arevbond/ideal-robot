package services

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/utils/e"
	"context"
	"time"
)

func (s *Service) GetHistories(limit int) ([]*models.History, error) {
	if limit < 0 {
		limit = 5
	}
	histories, err := s.db.GetHistories(context.Background(), limit)
	if err != nil {
		return nil, e.Wrap("can't get histories", err)
	}
	return histories, nil
}

func (s *Service) CreateHistory(text string) error {
	createHistory := &models.CreateHistory{Text: text, Type: 1, CreatedAt: time.Now()}
	err := s.db.CreateHistory(context.Background(), createHistory)
	if err != nil {
		return e.Wrap("can't create history", err)
	}
	return nil
}
