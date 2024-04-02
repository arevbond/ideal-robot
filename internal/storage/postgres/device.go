package postgres

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/utils/e"
	"context"
)

func (s *Storage) CreateDevice(ctx context.Context, device *models.CreateDevice) error {
	q := `INSERT INTO devices (room_id, name, type, status, write_topic, read_topic) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.ExecContext(ctx, q, device.RoomID, device.Name, device.Type, device.Status,
		device.WriteTopic, device.ReadTopic)
	if err != nil {
		return e.Wrap("cant create device in storage", err)
	}
	return nil
}

func (s *Storage) GetDeviceByID(ctx context.Context, id int) (*models.Device, error) {
	q := `SELECT * FROM devices WHERE id = $1`

	var device models.Device
	err := s.db.GetContext(ctx, &device, q, id)
	if err != nil {
		return nil, e.Wrap("can't get device by id from storage", err)
	}
	return &device, nil
}

func (s *Storage) GetDevicesByRoomID(ctx context.Context, roomID int, offset, limit int) ([]*models.Device, error) {
	q := `SELECT * FROM devices WHERE room_id = $1 OFFSET $2 LIMIT $3`

	devices := []*models.Device{}
	err := s.db.SelectContext(ctx, &devices, q, roomID, offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get devices by hub_id from storage", err)
	}
	return devices, nil
}

func (s *Storage) UpdateDevice(ctx context.Context, device *models.Device) error {
	q := `UPDATE devices SET room_id = $1, name = $2, type = $3, status = $4, write_topic = $5, read_topic = $6 WHERE id = $7`

	_, err := s.db.ExecContext(ctx, q, device.RoomID, device.Name, device.Type, device.Status,
		device.WriteTopic, device.ReadTopic, device.ID)
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
