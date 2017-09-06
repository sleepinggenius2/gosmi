package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/sleepinggenius2/gosmi/types"
)

type Import struct {
	Module string
	Name string
}

type Module struct {
	smiModule *C.struct_SmiModule
	ContactInfo string
	Description string
	Language types.Language
	Name string
	Organization string
	Path string
	Reference string
}

type Revision struct {
	Date syscall.Time_t
	Description string
}

func (m Module) GetIdentityNode() (node Node, ok bool) {
	smiIdentityNode := C.smiGetModuleIdentityNode(m.smiModule)
	if smiIdentityNode == nil {
		return
	}
	return CreateNode(smiIdentityNode), true
}

func (m Module) GetImports() (imports []Import) {
	for smiImport := C.smiGetFirstImport(m.smiModule); smiImport != nil; smiImport = C.smiGetNextImport(smiImport) {
		_import := Import{
			Module: C.GoString(smiImport.module),
			Name: C.GoString(smiImport.name),
		}
		imports = append(imports, _import)
	}
	return
}

func (m Module) GetNode(name string) (node Node, err error) {
	return GetNode(name, m)
}

func (m Module) GetNodes(kind ...types.NodeKind) (nodes []Node) {
	nodeKind := types.NodeAny
	if len(kind) > 0 && kind[0] != types.NodeUnknown {
		nodeKind = kind[0]
	}
	cNodeKind := C.SmiNodekind(nodeKind)
	for smiNode := C.smiGetFirstNode(m.smiModule, cNodeKind); smiNode != nil; smiNode = C.smiGetNextNode(smiNode, cNodeKind) {
		nodes = append(nodes, CreateNode(smiNode))
	}
	return
}

func (m Module) GetRevisions() (revisions []Revision) {
	for smiRevision := C.smiGetFirstRevision(m.smiModule); smiRevision != nil; smiRevision = C.smiGetNextRevision(smiRevision) {
		revision := Revision{
			Date: syscall.Time_t(smiRevision.date),
			Description: C.GoString(smiRevision.description),
		}
		revisions = append(revisions, revision)
	}
	return
}

func (m Module) GetTypes() (types []Type) {
	for smiType := C.smiGetFirstType(m.smiModule); smiType != nil; smiType = C.smiGetNextType(smiType) {
		types = append(types, CreateType(smiType))
	}
	return
}

func (m Module) GetRaw() (module *C.struct_SmiModule) {
	return m.smiModule
}

func (m *Module) SetRaw(smiModule *C.struct_SmiModule) {
	m.smiModule = smiModule
}

func CreateModule(smiModule *C.struct_SmiModule) (module Module) {
	module.SetRaw(smiModule)
	module.ContactInfo = C.GoString(smiModule.contactinfo)
	module.Description = C.GoString(smiModule.description)
	module.Language = types.Language(smiModule.language)
	module.Name = C.GoString(smiModule.name)
	module.Organization = C.GoString(smiModule.organization)
	module.Path = C.GoString(smiModule.path)
	module.Reference = C.GoString(smiModule.reference)
	return
}

func LoadModule(modulePath string) (moduleName string, err error) {
	cModulePath := C.CString(modulePath)
	defer C.free(unsafe.Pointer(cModulePath))

	cModuleName := C.smiLoadModule(cModulePath)
	if cModuleName == nil {
		err = fmt.Errorf("Could not load module at %s", modulePath)
		return
	}

	return C.GoString(cModuleName), nil
}

func GetLoadedModules() (modules []Module) {
	for smiModule := C.smiGetFirstModule(); smiModule != nil; smiModule = C.smiGetNextModule(smiModule) {
		modules = append(modules, CreateModule(smiModule))
	}
	return
}

func IsLoaded(moduleName string) bool {
	cModuleName := C.CString(moduleName)
	defer C.free(unsafe.Pointer(cModuleName))

	cStatus := C.smiIsLoaded(cModuleName)

	return C.int(cStatus) > 0
}

func GetModule(name string) (module Module, err error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	smiModule := C.smiGetModule(cName)
	if smiModule == nil {
		err = fmt.Errorf("Could not find module named %s", name)
		return
	}

	return CreateModule(smiModule), nil
}
