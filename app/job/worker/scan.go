package worker

import (
	"context"
	"fmt"
	"net"
)

type Scan struct {
	device deviceInformer
	port   int
}

type deviceInformer interface {
	GetThisDeviceID() string
}

func NewScan(device deviceInformer, port int) *Scan {
	return &Scan{
		device: device,
		port:   port,
	}
}

func (e *Scan) Init() error {
	return nil
}

func (e *Scan) Work(ctx context.Context) error {
	// broadcasting rather than unicasting
	broadcastAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", e.port))
	if err != nil {
		return err
	}
	udpConn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		return err
	}
	defer udpConn.Close()

	id := e.device.GetThisDeviceID()
	_, err = udpConn.Write([]byte(id))
	return err
}

func (e *Scan) Name() string {
	return "Scan"
}

func (e *Scan) Stop() {
}
