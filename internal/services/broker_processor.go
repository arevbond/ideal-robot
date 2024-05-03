package services

import (
	"HestiaHome/internal/clients/mqtt"
	"HestiaHome/internal/models"
	"HestiaHome/internal/storage"
	"HestiaHome/internal/utils/e"
	"context"
	"fmt"
	"log/slog"
	"time"
)

func processDataFromMQTT(log *slog.Logger, db storage.Storage, s *Service) {
	for {
		select {
		case data := <-mqtt.DevicesData:
			createDeviceData := extractDataByCategory(data)
			log.Info("receive data from broker", slog.Any("data", data))
			err := writeDataInStorage(s, log, db, data, createDeviceData)
			if err != nil {
				log.Error("can't write data in storage", slog.Any("error", err))
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func writeDataInStorage(s *Service, log *slog.Logger, db storage.Storage, deviceData *mqtt.DeviceData, createDeviceData *models.CreateDeviceData) error {
	_, err := db.GetDeviceByID(context.Background(), createDeviceData.DeviceID)
	if err == storage.ErrDeviceNotExist {
		err = db.CreateDevice(context.Background(), &models.Device{
			ID:       createDeviceData.DeviceID,
			Name:     deviceData.DeviceName,
			Category: deviceData.Category.Number(),
			Status:   true,
		})
		if err != nil {
			return e.Wrap("can't create device in storage", err)
		}
		err = s.CreateHistory(fmt.Sprintf(msgEventCreateDevice, deviceData.DeviceName))
		if err != nil {
			log.Error("can't create history", slog.Any("error", err))
		}
	} else if err != nil {
		return e.Wrap("can't get device by id", err)
	}

	err = db.CreateDeviceData(context.Background(), createDeviceData)
	if err != nil {
		return e.Wrap("can't create device data in storage", err)
	}
	log.Debug("success insert device data in storage")
	return nil
}

func extractDataByCategory(deviceData *mqtt.DeviceData) *models.CreateDeviceData {
	result := models.CreateDeviceData{DeviceID: deviceData.DeviceID, ReceivedAt: time.Now()}
	switch deviceData.Category {
	case mqtt.Temperature:
		if val, ok := deviceData.Data.(float64); ok {
			result.Value = models.TemperatureData{Value: val}
			result.Unit = "degrees"
		}
	case mqtt.Humidity:
		if val, ok := deviceData.Data.(float64); ok {
			result.Value = models.HumidityData{Value: int(val)}
			result.Unit = "%"
		}
	}
	return &result
}
