package postgres

import (
	"HestiaHome/internal/lib/e"
	"HestiaHome/internal/models"
	"context"
)

func (s *Storage) CreateDeviceData(ctx context.Context, deviceData *models.DeviceData) error {
	q := `INSERT INTO devices_data (device_id, value, unit, received_at) VALUES ($1, $2, $3, $4)`

	_, err := s.db.ExecContext(ctx, q, deviceData.DeviceID, deviceData.Value, deviceData.Unit, deviceData.ReceivedAt)
	if err != nil {
		return e.Wrap("can't create device data in storage", err)
	}
	return nil
}

func (s *Storage) GetDeviceDataByID(ctx context.Context, id int) (*models.DBDeviceData, error) {
	q := `SELECT * FROM devices_data WHERE id = $1`

	var deviceData models.DBDeviceData
	err := s.db.GetContext(ctx, &deviceData, q, id)
	if err != nil {
		return nil, e.Wrap("can't get device data by id from storage", err)
	}
	return &deviceData, nil
}

func (s *Storage) GetAllDeviceData(ctx context.Context, deviceID int, offset, limit int) ([]*models.DBDeviceData, error) {
	q := `SELECT * FROM devices_data WHERE device_id = $1 OFFSET $2 LIMIT $3`

	deviceData := []*models.DBDeviceData{}
	err := s.db.SelectContext(ctx, &deviceData, q, deviceID, offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get device data from storage", err)
	}
	return deviceData, nil
}

func (s *Storage) UpdateDeviceData(ctx context.Context, deviceData *models.DBDeviceData) error {
	q := `UPDATE devices_data SET device_id = $1, value = $2, unit = $3, received_at = $4 WHERE id = $5`

	_, err := s.db.ExecContext(ctx, q, deviceData.DeviceID, deviceData.Value, deviceData.Unit, deviceData.ReceivedAt, deviceData.ID)
	if err != nil {
		return e.Wrap("can't update device data in storage", err)
	}
	return nil
}

func (s *Storage) DeleteDeviceData(ctx context.Context, id int) error {
	q := `DELETE FROM devices_data WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		return e.Wrap("can't delete device data from storage", err)
	}
	return nil
}
