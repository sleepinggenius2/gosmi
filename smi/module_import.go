package smi

import (
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// SmiImport *smiGetFirstImport(SmiModule *smiModulePtr)
func GetFirstImport(smiModulePtr *types.SmiModule) *types.SmiImport {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	importPtr := modulePtr.Imports.First
	if importPtr == nil {
		return nil
	}
	return &importPtr.SmiImport
}

// SmiImport *smiGetNextImport(SmiImport *smiImportPtr)
func GetNextImport(smiImportPtr *types.SmiImport) *types.SmiImport {
	if smiImportPtr == nil {
		return nil
	}
	importPtr := (*internal.Import)(unsafe.Pointer(smiImportPtr))
	if importPtr.Next == nil {
		return nil
	}
	return &importPtr.Next.SmiImport
}

// int smiIsImported(SmiModule *smiModulePtr, SmiModule *importedModulePtr, char *importedName)
func IsImported(smiModulePtr *types.SmiModule, importedModulePtr *types.SmiModule, importedName string) bool {
	if smiModulePtr == nil || importedName == "" {
		return false
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	importPtr := modulePtr.Imports.Get(types.SmiIdentifier(importedName))
	if importPtr == nil {
		return false
	}
	return importedModulePtr == nil || importPtr.Module == importedModulePtr.Name
}
