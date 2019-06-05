package smi

import (
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// SmiType *smiGetType(SmiModule *smiModulePtr, char *type)
func GetType(smiModulePtr *types.SmiModule, typeName string) *types.SmiType {
	if typeName == "" {
		return nil
	}

	var modulePtr *internal.Module
	if smiModulePtr != nil {
		modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
		typePtr := modulePtr.Types.GetName(typeName)
		if typePtr == nil {
			return nil
		}
		return &typePtr.SmiType
	}
	for modulePtr = internal.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
		typePtr := modulePtr.Types.GetName(typeName)
		if typePtr != nil {
			return &typePtr.SmiType
		}
	}
	return nil
}

// SmiType *smiGetFirstType(SmiModule *smiModulePtr)
func GetFirstType(smiModulePtr *types.SmiModule) *types.SmiType {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	typePtr := modulePtr.Types.First
	if typePtr == nil {
		return nil
	}
	return &typePtr.SmiType
}

// SmiType *smiGetNextType(SmiType *smiTypePtr)
func GetNextType(smiTypePtr *types.SmiType) *types.SmiType {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Next == nil {
		return nil
	}
	return &typePtr.Next.SmiType
}

// SmiType *smiGetParentType(SmiType *smiTypePtr)
func GetParentType(smiTypePtr *types.SmiType) *types.SmiType {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Parent == nil {
		return nil
	}
	return &typePtr.Parent.SmiType
}

// SmiModule *smiGetTypeModule(SmiType *smiTypePtr)
func GetTypeModule(smiTypePtr *types.SmiType) *types.SmiModule {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.Module == nil {
		return nil
	}
	return &typePtr.Module.SmiModule
}

// int smiGetTypeLine(SmiType *smiTypePtr)
func GetTypeLine(smiTypePtr *types.SmiType) int {
	if smiTypePtr == nil {
		return 0
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	return typePtr.Line
}
