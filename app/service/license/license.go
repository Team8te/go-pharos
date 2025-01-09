package license

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	licenceFile = "license"
)

type Licenser struct {
	path string
	uuid uuid.UUID
}

func (l *Licenser) InitLicense() error {
	p := path.Join(l.path, licenceFile)
	_, err := os.Stat(p)
	if err == nil {
		l.loadLicence(p)
	}

	if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write([]byte(uuid.New().String()))
	if err != nil {
		return err
	}

	return f.Chmod(os.ModePerm)
}

func (l *Licenser) loadLicence(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}
	id, err := uuid.Parse(string(buf[:n]))
	if err != nil {
		return err
	}
	l.uuid = id
	return nil
}

func (l *Licenser) GetUUID() string {
	return l.uuid.String()
}

func NewLicenser(path string) (*Licenser, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}
	l := &Licenser{
		path: path,
	}

	return l, l.InitLicense()
}
