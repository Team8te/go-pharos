package device

import (
	"context"

	"github.com/go-pharos/app/ds"
)

type DeviceEndpoint struct {
	d deviceService
}

type deviceService interface {
	DeviceList(ctx context.Context, limit, offset int) ([]*ds.Device, error)
}

func NewDeviceEndpoint(d deviceService) *DeviceEndpoint {
	return &DeviceEndpoint{
		d: d,
	}
}
