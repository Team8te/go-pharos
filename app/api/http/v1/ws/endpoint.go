package ws

import (
	"context"

	"github.com/gorilla/websocket"
)

type store interface {
	MakeFile(ctx context.Context, name string, dataCh <-chan []byte) ([]byte, error)
}

type WSEndpoint struct {
	upgrader *websocket.Upgrader
	st       store
}

func NewWSEndpoint(st store) *WSEndpoint {
	return &WSEndpoint{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
		},
		st: st,
	}
}
