package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/sleepinggenius2/gosmi/models"
	"github.com/sleepinggenius2/gosmi/types"
)

type SmiModule struct {
	models.Module
	smiModule *C.struct_SmiModule
}

func (m SmiModule) GetIdentityNode() (node SmiNode, ok bool) {
	smiIdentityNode := C.smiGetModuleIdentityNode(m.smiModule)
	if smiIdentityNode == nil {
		return
	}
	return CreateNode(smiIdentityNode), true
}

func (m SmiModule) GetImports() (imports []models.Import) {
	for smiImport := C.smiGetFirstImport(m.smiModule); smiImport != nil; smiImport = C.smiGetNextImport(smiImport) {
		_import := models.Import{
			Module: C.GoString(smiImport.module),
			Name:   C.GoString(smiImport.name),
		}
		imports = append(imports, _import)
	}
	return
}

func (m SmiModule) GetNode(name string) (node SmiNode, err error) {
	return GetNode(name, m)
}

func (m SmiModule) GetNodeByOid(oid []uint) (node SmiNode, err error) {
	length := len(oid)
	var subid = make([]C.uint, length)
	defer C.free(unsafe.Pointer(subid))
	for i, o := range oid {
		subid[i] = C.uint(o)
	}
	ssp := (*C.SmiSubid)(unsafe.Pointer(&subid[0]))
	smiNode := C.smiGetNodeByOID((C.uint)(length), ssp)
	if smiNode == nil {
		fmt.Errorf("Could not find Oid %v in module %s", oid, m.Module.Name)
		return
	}
	return CreateNode(smiNode), nil
}

func (m SmiModule) GetNodes(kind ...types.NodeKind) (nodes []SmiNode) {
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

func (m SmiModule) GetRevisions() (revisions []models.Revision) {
	for smiRevision := C.smiGetFirstRevision(m.smiModule); smiRevision != nil; smiRevision = C.smiGetNextRevision(smiRevision) {
		revision := models.Revision{
			Date:        time.Unix(int64(smiRevision.date), 0),
			Description: C.GoString(smiRevision.description),
		}
		revisions = append(revisions, revision)
	}
	return
}

func (m SmiModule) GetType(name string) (outType SmiType, err error) {
	return GetType(name, m)
}

func (m SmiModule) GetTypes() (types []SmiType) {
	for smiType := C.smiGetFirstType(m.smiModule); smiType != nil; smiType = C.smiGetNextType(smiType) {
		types = append(types, CreateType(smiType))
	}
	return
}

func (m SmiModule) GetRaw() (module *C.struct_SmiModule) {
	return m.smiModule
}

func (m *SmiModule) SetRaw(smiModule *C.struct_SmiModule) {
	m.smiModule = smiModule
}

func CreateModule(smiModule *C.struct_SmiModule) (module SmiModule) {
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

func GetLoadedModules() (modules []SmiModule) {
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

func GetModule(name string) (module SmiModule, err error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	smiModule := C.smiGetModule(cName)
	if smiModule == nil {
		err = fmt.Errorf("Could not find module named %s", name)
		return
	}

	return CreateModule(smiModule), nil
}
