package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"unsafe"

	"github.com/sleepinggenius2/gosmi/types"
)

type Module struct {
	SmiModule *C.struct_SmiModule `json:"-"`
	ContactInfo string
	Description string
	Language types.Language
	Name string
	Organization string
	Path string
	Reference string
}

func (m Module) GetNodes(kind ...types.NodeKind) (nodes []Node) {
	nodeKind := types.NodeAny
	if len(kind) > 0 && kind[0] != types.NodeUnknown {
		nodeKind = kind[0]
	}
	cNodeKind := C.SmiNodekind(nodeKind)
	for smiNode := C.smiGetFirstNode(m.SmiModule, cNodeKind); smiNode != nil; smiNode = C.smiGetNextNode(smiNode, cNodeKind) {
		nodes = append(nodes, CreateNode(smiNode))
	}
	return
}

func (m Module) GetTypes() (types []Type) {
	for smiType := C.smiGetFirstType(m.SmiModule); smiType != nil; smiType = C.smiGetNextType(smiType) {
		types = append(types, CreateType(smiType))
	}
	return
}

func CreateModule(smiModule *C.struct_SmiModule) (module Module) {
	module.SmiModule = smiModule
	module.ContactInfo = C.GoString(smiModule.contactinfo)
	module.Description = C.GoString(smiModule.description)
	module.Language = types.Language(smiModule.language)
	module.Name = C.GoString(smiModule.name)
	module.Organization = C.GoString(smiModule.organization)
	module.Path = C.GoString(smiModule.path)
	module.Reference = C.GoString(smiModule.reference)
	return
}

func LoadModule(modulePath string) (moduleName string, ok bool) {
	cModulePath := C.CString(modulePath)
	defer C.free(unsafe.Pointer(cModulePath))

	cModuleName := C.smiLoadModule(cModulePath)
	if cModuleName == nil {
		return
	}

	return C.GoString(cModuleName), true
}

func GetLoadedModules() (modules []Module) {
	for smiModule := C.smiGetFirstModule(); smiModule != nil; smiModule = C.smiGetNextModule(smiModule) {
		modules = append(modules, CreateModule(smiModule))
	}
	return
}

func GetModule(name string) (module Module, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	smiModule := C.smiGetModule(cName)
	if smiModule == nil {
		return
	}

	return CreateModule(smiModule), true
}
