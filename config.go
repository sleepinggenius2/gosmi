package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
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
	newPath := C.CString(strings.Trim(path, ":"))
	defer C.free(unsafe.Pointer(newPath))
	C.smiSetPath(newPath)
}

func AppendPath(path string) {
	oldPath := GetPath()
	newPath := oldPath + ":" + path
	SetPath(newPath)
}

func PrependPath(path string) {
	oldPath := GetPath()
	newPath := path + ":" + oldPath
	SetPath(newPath)
}
