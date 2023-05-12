package udp

import (
	"context"

	"github.com/go-pharos/app/ds"
)

type deviceService interface {
	GetThisDeviceID() string
	AddDevice(ctx context.Context, d *ds.Device) error
}

type UDPEndpoint struct {
	service deviceService
}

func NewEndpoint(service deviceService) *UDPEndpoint {
	return &UDPEndpoint{
		service: service,
	}
}
