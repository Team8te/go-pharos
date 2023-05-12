package udp

import (
	"context"

	"github.com/go-pharos/app/ds"
	log "github.com/sirupsen/logrus"
)

func (e *UDPEndpoint) ReceiveHandler(ctx context.Context, message *ds.Device) error {
	if message.UUID == e.service.GetThisDeviceID() {
		return nil
	}

	log.Info("New device message: %v, from: %v", message.UUID, message.IP)
	return e.service.AddDevice(ctx, &ds.Device{
		UUID: message.UUID,
		IP:   message.IP,
	})
}
