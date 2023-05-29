package worker

import (
	"context"

	"github.com/go-pharos/app/ds"
)

type Mover struct {
	files  []string
	d      *ds.Device
	target string
}

func NewFileMover(d *ds.Device, files []string, target string) *Mover {
	return &Mover{
		files:  files,
		d:      d,
		target: target,
	}
}

func (m *Mover) Init() error {
	return nil
}

func (m *Mover) Work(ctx context.Context) error {

	return nil
}

func (m *Mover) Stop() {}

func (m *Mover) Name() string {
	return ""
}
