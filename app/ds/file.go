package ds

import "path/filepath"

type File struct {
	IsDir bool   `json:"is_dir"`
	Path  string `json:"path"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
}

func (f *File) Ext() string {
	return filepath.Ext(f.Name)
}
