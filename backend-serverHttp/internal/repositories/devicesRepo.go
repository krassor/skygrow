package repositories

import (
	"context"
	"errors"
	"github.com/krassor/skygrow/backend-serverHttp/internal/models/entities"
)

var (
	errDeviceAlreadyExist error = errors.New("device already exist in the database")
)

type DevicesRepository interface {
	FindAllDevices(ctx context.Context) ([]entities.Devices, error)
	CreateDevice(ctx context.Context, device entities.Devices) (entities.Devices, error)
	UpdateDevice(ctx context.Context, device entities.Devices) (entities.Devices, error)
	FindDeviceById(ctx context.Context, id uint) (entities.Devices, error)
}

func (r *repository) FindAllDevices(ctx context.Context) ([]entities.Devices, error) {
	var devices []entities.Devices
	tx := r.DB.WithContext(ctx).Find(&devices)
	if tx.Error != nil {
		return []entities.Devices{}, tx.Error
	}

	return devices, nil
}

func (r *repository) FindDeviceById(ctx context.Context, id uint) (entities.Devices, error) {
	var device entities.Devices
	tx := r.DB.WithContext(ctx).First(&device, id)
	if tx.Error != nil {
		return entities.Devices{}, tx.Error
	}
	return device, nil
}

func (r *repository) CreateDevice(ctx context.Context, device entities.Devices) (entities.Devices, error) {

	tx := r.DB.WithContext(ctx).Where(entities.Devices{DeviceIpAddress: device.DeviceIpAddress, DevicePort: device.DevicePort}).FirstOrCreate(&device)
	if tx.Error != nil {
		return entities.Devices{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entities.Devices{}, errDeviceAlreadyExist
	}
	return device, nil
}

func (r *repository) UpdateDevice(ctx context.Context, device entities.Devices) (entities.Devices, error) {
	tx := r.DB.WithContext(ctx).Save(&device)
	if tx.Error != nil {
		return entities.Devices{}, tx.Error
	}
	return device, nil
}
