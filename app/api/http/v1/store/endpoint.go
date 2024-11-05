package store

import (
	"context"
	"io"

	"github.com/Team8te/go-pharos/app/ds"
)

type StoreEndpoint struct {
	st storeService
}

type storeService interface {
	ListDir(ctx context.Context, path string) ([]*ds.File, error)
	MoveFiles(ctx context.Context, files []string, target *ds.Target) error
	UploadFile(ctx context.Context, reader io.Reader, name string) error
}

func NewStoreEndpoint(st storeService) *StoreEndpoint {
	return &StoreEndpoint{
		st: st,
	}
}
