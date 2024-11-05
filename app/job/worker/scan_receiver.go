package worker

import (
	"context"
	"fmt"
	"net"

	"github.com/Team8te/go-pharos/app/ds"
	log "github.com/sirupsen/logrus"
)

type ScanReceiver struct {
	port    int
	handler endpoint

	c net.PacketConn
}

type endpoint interface {
	ReceiveHandler(ctx context.Context, device *ds.Device) error
}

func NewScanReceiver(handler endpoint, port int) *ScanReceiver {
	return &ScanReceiver{
		port:    port,
		handler: handler,
	}
}

func (e *ScanReceiver) Init() error {
	var err error
	e.c, err = net.ListenPacket("udp", fmt.Sprintf(":%d", e.port))
	if err != nil {
		return err
	}

	return nil
}

func (e *ScanReceiver) Work(ctx context.Context) error {
	buf := make([]byte, 512)
	n, add, err := e.c.ReadFrom(buf)
	if err != nil {
		return err
	}

	go e.handle(ctx, &ds.Device{
		UUID: string(buf[:n]),
		IP:   add.String(),
	})

	return nil

}

func (e *ScanReceiver) handle(ctx context.Context, device *ds.Device) {
	defer func() {
		r := recover()
		if r != nil {
			log.WithField("worker", e.Name()).Errorf("Panic in process: %v", r)
		}
	}()
	err := e.handler.ReceiveHandler(ctx, device)
	if err != nil {
		log.WithField("worker", e.Name()).Errorf("Failed to process. Error: %v", err)
	}
}

func (e *ScanReceiver) Name() string {
	return "ScanReceiver"
}

func (e *ScanReceiver) Stop() {
	e.c.Close()
}
