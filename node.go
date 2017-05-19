package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	_ "fmt"
	"unsafe"

	"github.com/sleepinggenius2/gosmi/types"
)

type Node struct {
	SmiNode     *C.struct_SmiNode `json:"-"`
	Access      types.Access
	Decl        types.Decl
	Description string
	Kind        types.NodeKind
	Name        string
	Oid         []uint
	OidLen      int
	Status      types.Status
	Type        *Type
}

func (n Node) GetModule() (module Module) {
	smiModule := C.smiGetNodeModule(n.SmiNode)
	return CreateModule(smiModule)
}

func (n Node) GetSubtree() (nodes []Node) {
	first := true
	smiNode := n.SmiNode
	for oidlen := n.OidLen; smiNode != nil && (first || int(smiNode.oidlen) > oidlen); smiNode = C.smiGetNextNode(smiNode, C.SMI_NODEKIND_ANY) {
		node := CreateNode(smiNode)
		nodes = append(nodes, node)
		first = false
	}
	return
}

func (n Node) Render(flags types.Render) string {
	cRenderString := C.smiRenderNode(n.SmiNode, C.int(flags))

	return C.GoString(cRenderString)
}

func (n Node) RenderNumeric() string {
	cRenderString := C.smiRenderOID(n.SmiNode.oidlen, n.SmiNode.oid, C.int(types.RenderNumeric))

	return C.GoString(cRenderString)
}

func (n Node) RenderQualified() string {
	return n.Render(types.RenderQualified)
}

func CreateNode(smiNode *C.struct_SmiNode) (node Node) {
	node.SmiNode = smiNode
	node.Access = types.Access(smiNode.access)
	node.Decl = types.Decl(smiNode.decl)
	node.Description = C.GoString(smiNode.description)
	node.Kind = types.NodeKind(smiNode.nodekind)
	node.Name = C.GoString(smiNode.name)
	node.OidLen = int(smiNode.oidlen)
	node.Status = types.Status(smiNode.status)
	node.Type = CreateTypeFromNode(smiNode)

	length := node.OidLen
	subid := (*[1 << 30]C.SmiSubid)(unsafe.Pointer(smiNode.oid))[:length:length]

	node.Oid = make([]uint, length)
	for i := 0; i < length; i++ {
		node.Oid[i] = uint(subid[i])
	}
	return
}

func GetNode(name string, module ...Module) (node Node, ok bool) {
	var smiModule *C.struct_SmiModule
	if len(module) > 0 {
		smiModule = module[0].SmiModule
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	smiNode := C.smiGetNode(smiModule, cName)
	if smiNode == nil {
		return
	}
	return CreateNode(smiNode), true
}
