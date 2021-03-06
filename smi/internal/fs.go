package internal

import (
	"os"
	"path/filepath"
)

type NamedFS struct {
	Name string
	FS   FS
}

type pathFS string

func newPathFS(path string) NamedFS {
	return NamedFS{path, pathFS(path)}
}

func (p pathFS) Open(name string) (File, error) {
	filename := filepath.Join(string(p), name)
	return os.Open(filename)
}

func SetFS(fs ...NamedFS) {
	smiHandle.Paths = fs
}

func AppendFS(fs ...NamedFS) {
	smiHandle.Paths = append(smiHandle.Paths, fs...)
}

func PrependFS(fs ...NamedFS) {
	smiHandle.Paths = append(fs, smiHandle.Paths...)
}
