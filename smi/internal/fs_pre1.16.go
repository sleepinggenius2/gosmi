// +build !go1.16

package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File = io.ReadCloser

type DirEntry interface {
	Name() string
	IsDir() bool
}

type FS interface {
	Open(name string) (io.ReadCloser, error)
	ReadDir(name string) ([]DirEntry, error)
}

func (p pathFS) ReadDir(name string) ([]DirEntry, error) {
	path := string(p)
	if name != "." {
		path = filepath.Join(path, name)
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("Read directory entries: %w", err)
	}
	dirEntries := make([]DirEntry, 0, len(files))
	for _, file := range files {
		if file.IsDir() || file.Mode().IsRegular() {
			dirEntries = append(dirEntries, file)
		}
	}
	return dirEntries, nil
}
