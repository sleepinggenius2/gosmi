package smi

import (
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// SmiRefinement *smiGetFirstRefinement(SmiNode *smiComplianceNodePtr)
func GetFirstRefinement(smiComplianceNodePtr *types.SmiNode) *types.SmiRefinement {
	if smiComplianceNodePtr == nil {
		return nil
	}
	objPtr := (*internal.Object)(unsafe.Pointer(smiComplianceNodePtr))
	if objPtr.NodeKind != types.NodeCompliance || objPtr.RefinementList == nil || objPtr.RefinementList.Ptr == nil {
		return nil
	}
	return &objPtr.RefinementList.Ptr.(*internal.Refinement).SmiRefinement
}

// SmiRefinement *smiGetNextRefinement(SmiRefinement *smiRefinementPtr)
func GetNextRefinement(smiRefinementPtr *types.SmiRefinement) *types.SmiRefinement {
	if smiRefinementPtr == nil {
		return nil
	}
	refPtr := (*internal.Refinement)(unsafe.Pointer(smiRefinementPtr))
	if refPtr.List == nil || refPtr.List.Next == nil || refPtr.List.Next.Ptr == nil {
		return nil
	}
	return &refPtr.List.Next.Ptr.(*internal.Refinement).SmiRefinement
}

// SmiNode *smiGetRefinementNode(SmiRefinement *smiRefinementPtr)
func GetRefinementNode(smiRefinementPtr *types.SmiRefinement) *types.SmiNode {
	if smiRefinementPtr == nil {
		return nil
	}
	refinementPtr := (*internal.Refinement)(unsafe.Pointer(smiRefinementPtr))
	if refinementPtr.Object == nil {
		return nil
	}
	return refinementPtr.Object.GetSmiNode()
}

// SmiType *smiGetRefinementType(SmiRefinement *smiRefinementPtr)
func GetRefinementType(smiRefinementPtr *types.SmiRefinement) *types.SmiType {
	if smiRefinementPtr == nil {
		return nil
	}
	refinementPtr := (*internal.Refinement)(unsafe.Pointer(smiRefinementPtr))
	if refinementPtr.Type == nil || refinementPtr.Type.BaseType == types.BaseTypeUnknown {
		return nil
	}
	return &refinementPtr.Type.SmiType
}

// SmiType *smiGetRefinementWriteType(SmiRefinement *smiRefinementPtr)
func GetRefinementWriteType(smiRefinementPtr *types.SmiRefinement) *types.SmiType {
	if smiRefinementPtr == nil {
		return nil
	}
	refinementPtr := (*internal.Refinement)(unsafe.Pointer(smiRefinementPtr))
	if refinementPtr.WriteType == nil || refinementPtr.WriteType.BaseType == types.BaseTypeUnknown {
		return nil
	}
	return &refinementPtr.WriteType.SmiType
}

// int smiGetRefinementLine(SmiRefinement *smiRefinementPtr)
func GetRefinementLine(smiRefinementPtr *types.SmiRefinement) int {
	if smiRefinementPtr == nil {
		return 0
	}
	refinementPtr := (*internal.Refinement)(unsafe.Pointer(smiRefinementPtr))
	return refinementPtr.Line
}
