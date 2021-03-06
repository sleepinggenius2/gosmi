// +build go1.16

package internal

import (
	"io/fs"
	"os"
	"path/filepath"
)

type DirEntry = fs.DirEntry
type File = fs.File
type FS = fs.ReadDirFS

func (p pathFS) ReadDir(name string) ([]DirEntry, error) {
	path := string(p)
	if name != "." {
		path = filepath.Join(path, name)
	}
	return os.ReadDir(path)
}
