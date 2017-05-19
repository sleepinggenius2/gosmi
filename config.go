package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"os"
	"strings"
	"unsafe"
)

func Init() {
	ident := C.CString("gosmi")
	defer C.free(unsafe.Pointer(ident))

	initStatus := C.smiInit(ident)
	if int(initStatus) != 0 {
		panic("Failed to initialize")
	}
}

func Exit() {
	C.smiExit()
}

func GetPath() string {
	cPath := C.smiGetPath()
	return C.GoString(cPath)
}

func SetPath(path string) {
	newPath := C.CString(strings.Trim(path, string(os.PathListSeparator)))
	defer C.free(unsafe.Pointer(newPath))
	C.smiSetPath(newPath)
}

func AppendPath(path string) {
	oldPath := GetPath()
	newPath := oldPath + string(os.PathListSeparator) + path
	SetPath(newPath)
}

func PrependPath(path string) {
	oldPath := GetPath()
	newPath := path + string(os.PathListSeparator) + oldPath
	SetPath(newPath)
}

func ReadConfig(filename string, tag ...string) bool {
	configTag := "gosmi"
	if len(tag) > 0 {
		configTag = tag[0]
	}

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cTag := C.CString(configTag)
	defer C.free(unsafe.Pointer(cTag))

	cStatus := C.smiReadConfig(cFilename, cTag)

	return C.int(cStatus) == 0
}
