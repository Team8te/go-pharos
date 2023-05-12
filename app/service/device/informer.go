package device

import (
	"context"

	"github.com/go-pharos/app/ds"
)

type DeviceInformer struct {
	deviceUUID string
	repo       repository
}

type repository interface {
	AddDevice(ctx context.Context, d *ds.Device) error
}

func NewDeviceInformer(deviceUUID string, repo repository) *DeviceInformer {
	return &DeviceInformer{
		deviceUUID: deviceUUID,
		repo:       repo,
	}
}

func (i *DeviceInformer) GetThisDeviceID() string {
	return i.deviceUUID
}

func (i *DeviceInformer) AddDevice(ctx context.Context, d *ds.Device) error {
	return i.repo.AddDevice(ctx, d)
}

func (i *DeviceInformer) ListDevices(limit, offset int) ([]*ds.Device, error) {
	panic("not implemented")
}
