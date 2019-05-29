package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"fmt"
	"unsafe"

	"github.com/sleepinggenius2/gosmi/models"
	"github.com/sleepinggenius2/gosmi/types"
)

type SmiNode struct {
	models.Node
	smiNode *C.struct_SmiNode
	SmiType *SmiType
}

func (n SmiNode) GetModule() (module SmiModule) {
	smiModule := C.smiGetNodeModule(n.smiNode)
	return CreateModule(smiModule)
}

func (n SmiNode) GetSubtree() (nodes []SmiNode) {
	first := true
	smiNode := n.smiNode
	for oidlen := n.OidLen; smiNode != nil && (first || int(smiNode.oidlen) > oidlen); smiNode = C.smiGetNextNode(smiNode, C.SMI_NODEKIND_ANY) {
		node := CreateNode(smiNode)
		nodes = append(nodes, node)
		first = false
	}
	return
}

func (n SmiNode) Render(flags types.Render) string {
	cRenderString := C.smiRenderNode(n.smiNode, C.int(flags))

	return C.GoString(cRenderString)
}

func (n SmiNode) RenderNumeric() string {
	cRenderString := C.smiRenderOID(n.smiNode.oidlen, n.smiNode.oid, C.int(types.RenderNumeric))

	return C.GoString(cRenderString)
}

func (n SmiNode) RenderQualified() string {
	return n.Render(types.RenderQualified)
}

func (n SmiNode) GetRaw() (node *C.struct_SmiNode) {
	return n.smiNode
}

func (n *SmiNode) SetRaw(smiNode *C.struct_SmiNode) {
	n.smiNode = smiNode
}

func CreateNode(smiNode *C.struct_SmiNode) (node SmiNode) {
	node.SetRaw(smiNode)
	node.Access = types.Access(smiNode.access)
	node.Decl = types.Decl(smiNode.decl)
	node.Description = C.GoString(smiNode.description)
	node.Kind = types.NodeKind(smiNode.nodekind)
	node.Name = C.GoString(smiNode.name)
	node.OidLen = int(smiNode.oidlen)
	node.Status = types.Status(smiNode.status)
	node.SmiType = CreateTypeFromNode(smiNode)
	if node.SmiType != nil {
		node.Type = &node.SmiType.Type
	}

	length := node.OidLen
	subid := (*[1 << 30]C.SmiSubid)(unsafe.Pointer(smiNode.oid))[:length:length]

	node.Oid = make([]uint32, length)
	for i := 0; i < length; i++ {
		node.Oid[i] = uint32(subid[i])
	}
	return
}

func GetNode(name string, module ...SmiModule) (node SmiNode, err error) {
	var smiModule *C.struct_SmiModule
	if len(module) > 0 {
		smiModule = module[0].GetRaw()
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	smiNode := C.smiGetNode(smiModule, cName)
	if smiNode == nil {
		if len(module) > 0 {
			err = fmt.Errorf("Could not find node named %s in module %s", name, module[0].Name)
		} else {
			err = fmt.Errorf("Could not find node named %s", name)
		}
		return
	}
	return CreateNode(smiNode), nil
}
