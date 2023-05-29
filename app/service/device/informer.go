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
	ListDevice(ctx context.Context, limit, offset int) ([]*ds.Device, error)
	FindDevice(ctx context.Context, ip string) (*ds.Device, error)
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

func (i *DeviceInformer) DeviceList(ctx context.Context, limit, offset int) ([]*ds.Device, error) {
	return i.repo.ListDevice(ctx, limit, offset)
}

func (i *DeviceInformer) DeviceExists(ctx context.Context, ip string) (*ds.Device, error) {
	return i.repo.FindDevice(ctx, ip)
}
