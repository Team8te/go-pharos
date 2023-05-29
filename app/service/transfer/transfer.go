package transfer

import (
	"context"
	"sync"

	"github.com/go-pharos/app/ds"
)

var index int = 0

type taskMove struct {
	id    int
	t     *ds.Target
	files []string
}

type Transfer struct {
	mx    sync.RWMutex
	tasks map[int]context.CancelFunc
}

func NewTransfer() *Transfer {
	return &Transfer{}
}

func (t *Transfer) AddTask(ctx context.Context, target *ds.Target, files []string) (int, error) {
	bgCtx, cancel := context.WithCancel(context.Background())

	task := &taskMove{
		t:     target,
		files: files,
	}

	t.mx.Lock()
	defer t.mx.Unlock()

	result := index
	task.id = result
	t.tasks[result] = cancel
	index++

	go t.moveFiles(bgCtx, task)

	return result, nil
}

func (t *Transfer) moveFiles(ctx context.Context, tk *taskMove) {
}

func (t *Transfer) CancelTask(ctx context.Context, id int) error {
	t.mx.Lock()
	defer t.mx.Unlock()
	return nil
}
