package smi

import (
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// SmiRevision *smiGetFirstRevision(SmiModule *smiModulePtr)
func GetFirstRevision(smiModulePtr *types.SmiModule) *types.SmiRevision {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.FirstRevision == nil {
		return nil
	}
	return &modulePtr.FirstRevision.SmiRevision
}

// SmiRevision *smiGetNextRevision(SmiRevision *smiRevisionPtr)
func GetNextRevision(smiRevisionPtr *types.SmiRevision) *types.SmiRevision {
	if smiRevisionPtr == nil {
		return nil
	}
	revisionPtr := (*internal.Revision)(unsafe.Pointer(smiRevisionPtr))
	if revisionPtr.Next == nil {
		return nil
	}
	return &revisionPtr.Next.SmiRevision
}

// int smiGetRevisionLine(SmiRevision *smiRevisionPtr)
func GetRevisionLine(smiRevisionPtr *types.SmiRevision) int {
	if smiRevisionPtr == nil {
		return 0
	}
	revisionPtr := (*internal.Revision)(unsafe.Pointer(smiRevisionPtr))
	return revisionPtr.Line
}
