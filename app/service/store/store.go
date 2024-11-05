package store

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Team8te/go-pharos/app/ds"
	log "github.com/sirupsen/logrus"
)

type StoreKeeper struct {
	rootDir string
	repo    repository
	tr      transfer
}

type repository interface {
	FindDevice(ctx context.Context, ip string) (*ds.Device, error)
}

type transfer interface {
	AddTask(ctx context.Context, target *ds.Target, files []string) (int, error)
}

func NewStoreKeeper(root string, tr transfer) *StoreKeeper {
	root, err := filepath.Abs(root)
	if err != nil {
		panic(fmt.Errorf("root not exists: %v", root))
	}

	if fileInfo, err := os.Stat(root); err != nil || !fileInfo.IsDir() {
		panic(fmt.Errorf("root is not a dir: %v", root))
	}

	return &StoreKeeper{
		rootDir: root,
		tr:      tr,
	}
}

func (s *StoreKeeper) ListDir(ctx context.Context, path string) ([]*ds.File, error) {
	path, err := s.makePath(path)
	if err != nil {
		return nil, err
	}

	return s.ReadDir(ctx, path)
}

func (s *StoreKeeper) ReadDir(ctx context.Context, path string) ([]*ds.File, error) {
	content, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]*ds.File, 0, len(content))
	for _, f := range content {
		info, _ := f.Info()
		result = append(result, &ds.File{
			IsDir: f.IsDir(),
			Name:  f.Name(),
			Size:  info.Size(),
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

	/*
		_, err := s.repo.FindDevice(ctx, target.IP)
		if err != nil {
			return err
		}
	*/

	//s.tr.AddTask(ctx, target, files)

	return nil
}

func (s *StoreKeeper) FileExists(path string) error {
	if _, err := os.Stat(filepath.Join(s.rootDir, path)); err == nil {
		return nil
	} else {
		return err
	}
}

func (s *StoreKeeper) MakeFile(ctx context.Context, name string, dataCh <-chan []byte) ([]byte, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case ch := <-dataCh:
			if len(ch) == 0 {
				return ds.Hash(file)
			}
			file.Write(ch)
		}
	}
}

func (s *StoreKeeper) UploadFile(ctx context.Context, reader io.Reader, name string) error {
	path, err := s.makePath(name)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	size, err := io.Copy(file, reader)
	log.WithFields(log.Fields{
		"filepath": path,
		"size":     size,
	}).Debugf("UploadFile")
	return err
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
