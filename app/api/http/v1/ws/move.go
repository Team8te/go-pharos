package ws

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const pongWait = 60 * time.Second

func (e *WSEndpoint) Move(c echo.Context) error {
	log.WithField("api", "http").Debugf("Move")
	conn, err := e.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	defer conn.Close()
	ctx := c.Request().Context()

	state := fileState{
		name:    "",
		fileCh:  make(chan []byte, 1024),
		hashSum: make(chan []byte),
	}
	defer close(state.fileCh)
	defer close(state.hashSum)

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		log.WithFields(log.Fields{
			"api":  "websocket",
			"size": len(msg),
		}).Debugf("Move")

		err = e.processChunk(ctx, &state, t, msg)
		if err != nil {
			break
		}
	}

	h := <-state.hashSum

	oh, err := e.readHash(ctx, conn)
	if err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			return err
		}
	}

	if bytes.Equal(h, oh) {
		log.WithFields(log.Fields{
			"local hash":  h,
			"origin hash": oh,
		}).Errorf("Invalid origin hash")
	}

	log.WithFields(log.Fields{
		"api":         "websocket",
		"local hash":  h,
		"origin hash": oh,
	}).Debugf("Move")
	return state.err
}

type fileState struct {
	name   string
	fileCh chan []byte

	hashSum chan []byte
	err     error
}

func (e *WSEndpoint) processChunk(ctx context.Context, state *fileState, t int, msg []byte) error {
	switch t {
	case websocket.BinaryMessage:
		state.fileCh <- msg
	case websocket.TextMessage:
		state.name = string(msg)
		go func() {
			h, err := e.st.MakeFile(ctx, state.name, state.fileCh)
			state.err = err
			state.hashSum <- h
		}()
	default:
		return fmt.Errorf("unknown chunk type")
	}

	return nil
}

func (e *WSEndpoint) readHash(ctx context.Context, conn *websocket.Conn) ([]byte, error) {
	t, msg, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	if t != websocket.BinaryMessage {
		return nil, fmt.Errorf("unknown chunk type")
	}

	return msg, nil
}
