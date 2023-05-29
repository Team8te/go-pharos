package store

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-pharos/app/ds"
)

type StoreKeeper struct {
	rootDir string
	repo    repository
}

type repository interface {
	FindDevice(ctx context.Context, ip string) (*ds.Device, error)
}

func NewStoreKeeper(root string) *StoreKeeper {
	root, err := filepath.Abs(root)
	if err != nil {
		panic(fmt.Errorf("root not exists: %v", root))
	}

	if fileInfo, err := os.Stat(root); err != nil || !fileInfo.IsDir() {
		panic(fmt.Errorf("root is not a dir: %v", root))
	}

	return &StoreKeeper{
		rootDir: root,
	}
}

func (s *StoreKeeper) ListDir(ctx context.Context, path string) ([]*ds.File, error) {
	path, err := s.makePath(path)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]*ds.File, 0, len(content))
	for _, f := range content {
		result = append(result, &ds.File{
			IsDir: f.IsDir(),
			Name:  f.Name(),
			Size:  f.Size(),
			Path:  path[len(s.rootDir):],
		})
	}

	return result, nil
}

func (s *StoreKeeper) MoveFiles(ctx context.Context, files []string, target *ds.Target) error {
	for _, f := range files {
		err := s.FileExists(f)
		if err != nil {
			return err
		}
	}

	_, err := s.repo.FindDevice(ctx, target.IP)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreKeeper) FileExists(path string) error {
	if _, err := os.Stat(filepath.Join(s.rootDir, path)); err == nil {
		return nil
	} else {
		return err
	}
}

func (s *StoreKeeper) makePath(path string) (string, error) {
	path, err := filepath.Abs(filepath.Join(s.rootDir, path))
	if err != nil {
		return "", err
	}

	i := strings.Index(path, s.rootDir)
	if i != 0 {
		return "", fmt.Errorf("invalid path")
	}

	return path, nil
}
