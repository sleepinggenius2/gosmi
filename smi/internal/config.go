package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
	names := make([]string, len(smiHandle.Paths))
	for i, fs := range smiHandle.Paths {
		names[i] = fs.Name
	}
	return strings.Join(names, string(os.PathListSeparator))
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
		path = filepath.Join(homedir, path[1:])
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("Get absolute path for '%s': %w", path, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("Cannot stat '%s': %w", path, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("'%s' is not a directory", path)
	}
	return path, nil
}

func SetPath(path ...string) {
	pathLen := len(path)
	if pathLen == 0 {
		return
	}
	if path[0] == "" {
		appendPath(path[1:]...)
	} else if path[pathLen-1] == "" {
		prependPath(path[:pathLen-1]...)
	} else {
		smiHandle.Paths = make([]NamedFS, 0, pathLen)
		for _, p := range path {
			if p, err := expandPath(p); err == nil {
				smiHandle.Paths = append(smiHandle.Paths, newPathFS(p))
			}
		}
	}
}

func appendPath(path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]NamedFS, len(smiHandle.Paths), len(smiHandle.Paths)+len(path))
	copy(paths, smiHandle.Paths)
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, newPathFS(p))
		}
	}
	smiHandle.Paths = paths
}

func prependPath(path ...string) {
	if len(path) == 0 {
		return
	}
	paths := make([]NamedFS, 0, len(smiHandle.Paths)+len(path))
	for _, p := range path {
		if p, err := expandPath(p); err == nil {
			paths = append(paths, newPathFS(p))
		}
	}
	paths = append(paths, smiHandle.Paths...)
	smiHandle.Paths = paths
}
