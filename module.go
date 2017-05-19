package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"syscall"
	"unsafe"

	"github.com/sleepinggenius2/gosmi/types"
)

type Import struct {
	Module string
	Name string
}

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

type Revision struct {
	Date syscall.Time_t
	Description string
}

func (m Module) GetIdentityNode() (node Node, ok bool) {
	smiIdentityNode := C.smiGetModuleIdentityNode(m.SmiModule)
	if smiIdentityNode == nil {
		return
	}
	return CreateNode(smiIdentityNode), true
}

func (m Module) GetImports() (imports []Import) {
	for smiImport := C.smiGetFirstImport(m.SmiModule); smiImport != nil; smiImport = C.smiGetNextImport(smiImport) {
		_import := Import{
			Module: C.GoString(smiImport.module),
			Name: C.GoString(smiImport.name),
		}
		imports = append(imports, _import)
	}
	return
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

func (m Module) GetRevisions() (revisions []Revision) {
	for smiRevision := C.smiGetFirstRevision(m.SmiModule); smiRevision != nil; smiRevision = C.smiGetNextRevision(smiRevision) {
		revision := Revision{
			Date: syscall.Time_t(smiRevision.date),
			Description: C.GoString(smiRevision.description),
		}
		revisions = append(revisions, revision)
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
