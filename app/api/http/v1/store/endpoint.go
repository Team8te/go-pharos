package store

import (
	"context"

	"github.com/go-pharos/app/ds"
)

type StoreEndpoint struct {
	st storeService
}

type storeService interface {
	ListDir(ctx context.Context, path string) ([]*ds.File, error)
	MoveFiles(ctx context.Context, files []string, target *ds.Target) error
}

func NewStoreEndpoint(st storeService) *StoreEndpoint {
	return &StoreEndpoint{
		st: st,
	}
}
