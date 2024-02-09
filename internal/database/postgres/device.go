package postgres

import (
	"HestiaHome/internal/lib/e"
	"HestiaHome/internal/models"
	"context"
)

func (s *Storage) CreateDevice(ctx context.Context, device *models.Device) error {
	q := `INSERT INTO devices (hub_id, name, type, location, status) VALUES ($1, $2, $3, $4, $5)`

	_, err := s.db.ExecContext(ctx, q, device.HubID, device.Name, device.Type, device.Location, device.Status)
	if err != nil {
		return e.Wrap("cant create device in storage", err)
	}
	return nil
}

func (s *Storage) GetDeviceByID(ctx context.Context, id int) (*models.DBDevice, error) {
	q := `SELECT * FROM devices WHERE id = $1`

	var device models.DBDevice
	err := s.db.GetContext(ctx, &device, q, id)
	if err != nil {
		return nil, e.Wrap("can't get device by id from storage", err)
	}
	return &device, nil
}

func (s *Storage) GetDevicesByHubID(ctx context.Context, hubID int, offset, limit int) ([]*models.DBDevice, error) {
	q := `SELECT * FROM devices WHERE hub_id = $1 OFFSET $2 LIMIT $3`

	devices := []*models.DBDevice{}
	err := s.db.SelectContext(ctx, &devices, q, hubID, offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get devices by hub_id from storage", err)
	}
	return devices, nil
}

func (s *Storage) UpdateDevice(ctx context.Context, device *models.DBDevice) error {
	q := `UPDATE devices SET hub_id = $1, name = $2, type = $3, location = $4, status = $5 WHERE id = $6`

	_, err := s.db.ExecContext(ctx, q, device.HubID, device.Name, device.Type,
		device.Location, device.Status, device.ID)
	if err != nil {
		return e.Wrap("can't update device in storage", err)
	}
	return nil
}

func (s *Storage) DeleteDevice(ctx context.Context, id int) error {
	q := `DELETE FROM devices WHERE id = $1`
	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return e.Wrap("can't delete device from storage", err)
	}
	return nil
}
