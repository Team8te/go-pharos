package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Team8te/go-pharos/app/ds"
	"github.com/Team8te/go-pharos/app/service/transfer"
	"github.com/gorilla/websocket"
)

// NOTE: for simplicity, error check is omitted
func uploadLargeFile(addr, filePath string, chunkSize int) error {
	//open file and retrieve info
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	u := url.URL{Scheme: "ws", Host: addr, Path: "/v1/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	defer c.Close()

	//write file
	buf := make([]byte, chunkSize)
	for {
		_, err := file.Read(buf)
		if err != nil {
			break
		}
		err = c.WriteMessage(websocket.BinaryMessage, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendContent(u *url.URL, path string) error {
	ctx := context.Background()

	t := transfer.NewTransfer(1024)
	if !checkFileExists(path) {
		return fmt.Errorf("file not exists: %v", path)
	}

	path, _ = filepath.Abs(path)

	// "/Users/aleksasorokin/Downloads/Fork-2.27.dmg",

	return t.MoveFiles(ctx, &ds.TaskMove{
		URL:   u,
		Files: []string{path},
	})
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
