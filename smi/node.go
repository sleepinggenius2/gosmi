package smi

import (
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// SmiNode *smiGetNode(SmiModule *smiModulePtr, const char *name)
func GetNode(smiModulePtr *types.SmiModule, name string) *types.SmiNode {
	if name == "" {
		return nil
	}
	var modulePtr *internal.Module
	if smiModulePtr != nil {
		modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
		objPtr := modulePtr.Objects.GetName(name)
		if objPtr == nil {
			return nil
		}
		return objPtr.GetSmiNode()
	}
	for modulePtr = internal.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
		objPtr := modulePtr.Objects.GetName(name)
		if objPtr != nil {
			return objPtr.GetSmiNode()
		}
	}
	return nil
}

// SmiNode *smiGetNodeByOID(unsigned int oidlen, SmiSubid oid[])
func GetNodeByOID(oid types.Oid) *types.SmiNode {
	if len(oid) == 0 || internal.Root() == nil {
		return nil
	}
	var parentPtr, nodePtr *internal.Node = nil, internal.Root()
	for i := 0; i < len(oid) && nodePtr != nil; i++ {
		parentPtr, nodePtr = nodePtr, nodePtr.Children.Get(oid[i])
	}
	if nodePtr == nil {
		nodePtr = parentPtr
	}
	if nodePtr == nil || nodePtr.FirstObject == nil {
		return nil
	}
	return nodePtr.FirstObject.GetSmiNode()
}

// SmiNode *smiGetFirstNode(SmiModule *smiModulePtr, SmiNodekind nodekind)
func GetFirstNode(smiModulePtr *types.SmiModule, nodekind types.NodeKind) *types.SmiNode {
	if smiModulePtr == nil {
		return nil
	}
	var (
		modulePtr *internal.Module
		nodePtr   *internal.Node
		objPtr    *internal.Object
	)
	modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.PrefixNode != nil {
		nodePtr = modulePtr.PrefixNode
	} else if internal.Root() != nil {
		nodePtr = internal.Root().Children.First
	}
	for nodePtr != nil {
		objPtr = internal.GetNextChildObject(nodePtr, modulePtr, nodekind)
		if objPtr != nil {
			return objPtr.GetSmiNode()
		}
		if nodePtr.Children.First != nil {
			nodePtr = nodePtr.Children.First
		} else if nodePtr.Next != nil {
			nodePtr = nodePtr.Next
		} else {
			if nodePtr.Parent == nil {
				return nil
			}
			for nodePtr.Parent != nil && nodePtr.Next == nil {
				nodePtr = nodePtr.Parent
			}
			nodePtr = nodePtr.Next
		}
	}
	return nil
}

// SmiNode *smiGetNextNode(SmiNode *smiNodePtr, SmiNodekind nodekind)
func GetNextNode(smiNodePtr *types.SmiNode, nodekind types.NodeKind) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Module == nil || objPtr.Node == nil {
		return nil
	}
	nodePtr := objPtr.Node
	modulePtr := objPtr.Module
	for nodePtr != nil {
		if nodePtr.Children.First != nil {
			nodePtr = nodePtr.Children.First
		} else if nodePtr.Next != nil {
			nodePtr = nodePtr.Next
		} else {
			for nodePtr.Parent != nil && nodePtr.Next == nil {
				nodePtr = nodePtr.Parent
			}
			nodePtr = nodePtr.Next
			if nodePtr == nil || !nodePtr.Oid.ChildOf(modulePtr.PrefixNode.Oid) {
				return nil
			}
		}
		objPtr = internal.GetNextChildObject(nodePtr, modulePtr, nodekind)
		if objPtr != nil {
			return objPtr.GetSmiNode()
		}
	}
	return nil
}

// SmiNode *smiGetParentNode(SmiNode *smiNodePtr)
func GetParentNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Node == nil || objPtr.Node.Parent == nil || objPtr.Node.Flags.Has(internal.FlagRoot) {
		return nil
	}
	var parentPtr *internal.Object
	if objPtr.Module != nil {
		parentPtr = internal.FindObjectByModuleAndNode(objPtr.Module, objPtr.Node)
		if parentPtr != nil {
			importPtr := objPtr.Module.Imports.Get(parentPtr.Name)
			if importPtr != nil {
				parentPtr = internal.FindObjectByModuleNameAndNode(string(importPtr.Module), objPtr.Node)
			} else {
				parentPtr = nil
			}
		}
	}
	if parentPtr == nil {
		parentPtr = internal.FindObjectByNode(objPtr.Node)
	}
	if parentPtr == nil {
		return nil
	}
	return parentPtr.GetSmiNode()
}

// SmiNode *smiGetRelatedNode(SmiNode *smiNodePtr)
func GetRelatedNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Related == nil {
		return nil
	}
	return objPtr.Related.GetSmiNode()
}

// SmiNode *smiGetFirstChildNode(SmiNode *smiNodePtr)
func GetFirstChildNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Node == nil || objPtr.Node.Children.First == nil {
		return nil
	}
	nodePtr := objPtr.Node.Children.First
	objPtr = internal.FindObjectByModuleAndNode(objPtr.Module, nodePtr)
	if objPtr == nil {
		objPtr = internal.FindObjectByNode(nodePtr)
	}
	if objPtr == nil {
		return nil
	}
	return objPtr.GetSmiNode()
}

// SmiNode *smiGetNextChildNode(SmiNode *smiNodePtr)
func GetNextChildNode(smiNodePtr *types.SmiNode) *types.SmiNode {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Node == nil || objPtr.Node.Next == nil {
		return nil
	}
	nodePtr := objPtr.Node.Next
	objPtr = internal.FindObjectByModuleAndNode(objPtr.Module, nodePtr)
	if objPtr == nil {
		objPtr = internal.FindObjectByNode(nodePtr)
	}
	if objPtr == nil {
		return nil
	}
	return objPtr.GetSmiNode()
}

// SmiModule *smiGetNodeModule(SmiNode *smiNodePtr)
func GetNodeModule(smiNodePtr *types.SmiNode) *types.SmiModule {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Module == nil {
		return nil
	}
	return &objPtr.Module.SmiModule
}

// SmiType *smiGetNodeType(SmiNode *smiNodePtr)
func GetNodeType(smiNodePtr *types.SmiNode) *types.SmiType {
	if smiNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	if objPtr.Type == nil {
		return nil
	}
	return &objPtr.Type.SmiType
}

// int smiGetNodeLine(SmiNode *smiNodePtr)
func GetNodeLine(smiNodePtr *types.SmiNode) int {
	if smiNodePtr == nil {
		return 0
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiNodePtr))
	return objPtr.Line
}
