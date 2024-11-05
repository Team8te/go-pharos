package transfer

import (
	"context"
	"io"
	"net/url"
	"os"
	"sync"

	"github.com/Team8te/go-pharos/app/ds"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var index int = 0

type Transfer struct {
	mx    sync.RWMutex
	tasks map[int]context.CancelFunc

	chunkSize int
}

func NewTransfer(chunkSize int) *Transfer {
	return &Transfer{
		tasks:     make(map[int]context.CancelFunc),
		chunkSize: chunkSize,
	}
}

func (t *Transfer) AddTask(ctx context.Context, target *ds.Target, files []string) (int, error) {
	bgCtx, cancel := context.WithCancel(context.Background())

	task := &ds.TaskMove{
		URL:   &url.URL{Scheme: "ws", Host: target.IP, Path: "/v1/ws"},
		Files: files,
	}

	t.mx.Lock()
	defer t.mx.Unlock()

	result := index
	task.ID = result
	t.tasks[result] = cancel
	index++

	go func() {
		err := t.MoveFiles(bgCtx, task)
		log.WithField("service", "AddTask").
			WithError(err).
			Error()
	}()

	return result, nil
}

func (t *Transfer) MoveFiles(ctx context.Context, task *ds.TaskMove) error {
	log.Infof("connecting to %s", task.URL.String())

	c, _, err := websocket.DefaultDialer.Dial(task.URL.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	for _, p := range task.Files {
		err := t.sendFile(ctx, c, p)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithField("service", "MoveFiles").
					WithError(err).
					Error()
				return err
			}
		}
	}

	return nil
}

func (t *Transfer) sendFile(ctx context.Context, conn *websocket.Conn, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	err = t.sendFileName(ctx, conn, fi.Name())
	if err != nil {
		return err
	}

	defer file.Close()

	buf := make([]byte, t.chunkSize)
	ss := int64(0)
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, err := file.Read(buf)
			if err == io.EOF {
				return t.sendHash(ctx, conn, file)
			}
			if err != nil {
				return err
			}
			err = conn.WriteMessage(websocket.BinaryMessage, buf)
			if err != nil {
				return err
			}

			ss += int64(len(buf))

			log.WithFields(log.Fields{
				"file":       path,
				"total size": fi.Size(),
				"send size":  ss,
				"chunk size": len(buf),
			}).Debug("sendFile")

			if len(buf) < n {
				return nil
			}
		}
	}
}

func (t *Transfer) sendFileName(ctx context.Context, conn *websocket.Conn, name string) error {
	log.WithField("file", name).Debug("sendFileName")
	return conn.WriteMessage(websocket.TextMessage, []byte(name))
}

func (t *Transfer) sendHash(ctx context.Context, conn *websocket.Conn, file *os.File) error {
	err := conn.WriteMessage(websocket.BinaryMessage, []byte{})
	if err != nil {
		return err
	}
	h, err := ds.Hash(file)
	if err != nil {
		return err
	}
	log.WithField("hash", h).Debug("sendHash")
	return conn.WriteMessage(websocket.BinaryMessage, h)
}

func (t *Transfer) CancelTask(ctx context.Context, id int) error {
	t.mx.Lock()
	defer t.mx.Unlock()
	return nil
}
