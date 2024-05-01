package postgres

import (
	"HestiaHome/internal/models"
	"HestiaHome/internal/storage"
	"HestiaHome/internal/utils/e"
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

func (s *Storage) GetDevicesWithData(ctx context.Context) ([]*models.DeviceWithData, error) {
	q := `SELECT d.id AS device_id, d.name AS device_name, d.category as category, dd.value, dd.unit, dd.received_at
		FROM devices d
		JOIN (
			SELECT device_id, MAX(received_at) AS max_received_at
			FROM devices_data
			GROUP BY device_id
		) max_dd ON d.id = max_dd.device_id
		JOIN devices_data dd ON d.id = dd.device_id AND max_dd.max_received_at = dd.received_at;
`
	devices := []*DeviceWithDataEntity{}
	err := s.db.SelectContext(ctx, &devices, q)
	if err != nil {
		s.log.Error("can't get devices with data", "error", err)
		return nil, e.Wrap("can't get devices with data", err)
	}

	result := []*models.DeviceWithData{}
	for _, deviceFromDB := range devices {
		device, err := deviceFromDB.convertToModel()
		if err != nil {
			s.log.Error("can't convert db entity to model", "error", err)
			continue
		} else {
			result = append(result, device)
		}
	}
	return result, nil
}

func (s *Storage) GetDevicesWithDataByID(ctx context.Context, id int) ([]*models.DeviceWithData, error) {
	q := `SELECT d.id AS device_id, d.name AS device_name, d.category as category, dd.value, dd.unit, dd.received_at
		FROM devices d
		JOIN (
			SELECT device_id, MAX(received_at) AS max_received_at
			FROM devices_data
			GROUP BY device_id
		) max_dd ON d.id = max_dd.device_id
		JOIN devices_data dd ON d.id = dd.device_id AND max_dd.max_received_at = dd.received_at
		WHERE d.room_id = $1;
`
	devices := []*DeviceWithDataEntity{}
	err := s.db.SelectContext(ctx, &devices, q, id)
	if err != nil {
		s.log.Error("can't get devices with data", "error", err)
		return nil, e.Wrap("can't get devices with data", err)
	}

	result := []*models.DeviceWithData{}
	for _, deviceFromDB := range devices {
		device, err := deviceFromDB.convertToModel()
		if err != nil {
			s.log.Error("can't convert db entity to model", "error", err)
			continue
		} else {
			result = append(result, device)
		}
	}
	return result, nil
}

func (s *Storage) GetDevices(ctx context.Context) ([]*models.Device, error) {
	q := `SELECT * FROM devices`

	devices := []*models.Device{}
	err := s.db.SelectContext(ctx, &devices, q)
	if err != nil {
		return nil, e.Wrap("can't get device from storage", err)
	}
	return devices, nil
}

func (s *Storage) CreateDevice(ctx context.Context, device *models.Device) error {
	var q string
	if device.RoomID == nil {
		// Если room_id равен 0, пропустить room_id в запросе
		q = `INSERT INTO devices (id, name, category, status, hidden) VALUES ($1, $2, $3, $4, $5)`
		_, err := s.db.ExecContext(ctx, q, device.ID, device.Name, device.Category, device.Status, device.Hidden)
		if err != nil {
			return e.Wrap("cant create device in storage", err)
		}
		return nil
	}

	// Если room_id не равен 0, включить room_id в запрос
	q = `INSERT INTO devices (id, room_id, name, category, status, hidden) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := s.db.ExecContext(ctx, q, device.ID, device.RoomID, device.Name, device.Category, device.Status, device.Hidden)
	if err != nil {
		return e.Wrap("cant create device in storage", err)
	}
	return nil
}

func (s *Storage) GetDeviceByID(ctx context.Context, id int) (*models.Device, error) {
	q := `SELECT * FROM devices WHERE id = $1`

	var device models.Device
	err := s.db.GetContext(ctx, &device, q, id)

	if err == sql.ErrNoRows {
		return nil, storage.ErrDeviceNotExist
	}

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
	q := `UPDATE devices SET room_id = $1, name = $2, category = $3, hidden = $4, status = $5  WHERE id = $6`

	_, err := s.db.ExecContext(ctx, q, device.RoomID, device.Name, device.Category, device.Hidden, device.Status, device.ID)
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

type DeviceWithDataEntity struct {
	ID         int       `db:"device_id"`
	Name       string    `db:"device_name"`
	Category   int       `db:"category"`
	Value      []byte    `db:"value"`
	Unit       string    `db:"unit"`
	ReceivedAt time.Time `db:"received_at"`
}

func (d *DeviceWithDataEntity) convertToModel() (*models.DeviceWithData, error) {
	var value models.Value
	err := json.Unmarshal(d.Value, &value)
	if err != nil {
		return nil, e.Wrap("can't convert value from db to json", err)
	}
	return &models.DeviceWithData{
		ID:         d.ID,
		Name:       d.Name,
		Category:   d.Category,
		Value:      value,
		Unit:       d.Unit,
		ReceivedAt: d.ReceivedAt,
	}, nil
}
