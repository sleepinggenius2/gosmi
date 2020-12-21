package internal

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var moduleFiles *ModuleFiles

func init() {
	moduleFiles = &ModuleFiles{
		m: make(map[string][]byte),
	}
}

func Init(handleName string) bool {
	smiHandle = findHandleByName(handleName)
	if smiHandle != nil {
		return true
	}
	smiHandle = addHandle(handleName)
	return initData()
}

func Exit() {
	if smiHandle == nil {
		return
	}
	freeData()
	removeHandle(smiHandle)
}

func GetPath() string {
	return strings.Join(smiHandle.Paths, string(os.PathListSeparator))
}

func expandPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("Path is empty")
	}
	if path[0] == '~' {
		homedir, err := os.UserHomeDir()
		if err != nil {
			return "", errors.New("Cannot expand homedir")
		}
		path = filepath.Join(homedir, path)
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", errors.Wrapf(err, "Get absolute path for '%s'", path)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", errors.Wrapf(err, "Cannot stat '%s'", path)
	}
	if !info.IsDir() {
		return "", errors.Errorf("'%s' is not a directory", path)
	}
	return path, nil
}

func SetPath(path ...string) {
	pathLen := len(path)
	if pathLen == 0 {
		return
	}
	if path[0] == "" {
		AppendPath(path[1:]...)
	} else if path[pathLen-1] == "" {
		PrependPath(path[:pathLen-1]...)
	} else {
		smiHandle.Paths = make([]string, 0, pathLen)
		for _, p := range path {
			if p, err := expandPath(p); err == nil {
				smiHandle.Paths = append(smiHandle.Paths, p)
			}
		}
	}
}

func AppendPath(path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]string, len(smiHandle.Paths), len(smiHandle.Paths)+len(path))
	copy(paths, smiHandle.Paths)
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, p)
		}
	}
	smiHandle.Paths = paths
}

func PrependPath(path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]string, 0, len(smiHandle.Paths)+len(path))
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, p)
		}
	}
	paths = append(paths, smiHandle.Paths...)
	smiHandle.Paths = paths
}

func AddModuleFile(name string, data []byte) {
	moduleFiles.Add(name, data)
}
