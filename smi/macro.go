package smi

import (
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// SmiMacro *smiGetMacro(SmiModule *smiModulePtr, char *macro)
func GetMacro(smiModulePtr *types.SmiModule, macro string) *types.SmiMacro {
	if macro == "" {
		return nil
	}

	var modulePtr *internal.Module
	if smiModulePtr != nil {
		modulePtr = (*internal.Module)(unsafe.Pointer(smiModulePtr))
		macroPtr := modulePtr.Macros.GetName(macro)
		if macroPtr == nil {
			return nil
		}
		return &macroPtr.SmiMacro
	}
	for modulePtr = internal.GetFirstModule(); modulePtr != nil; modulePtr = modulePtr.Next {
		macroPtr := modulePtr.Macros.GetName(macro)
		if macroPtr != nil {
			return &macroPtr.SmiMacro
		}
	}
	return nil
}

// SmiMacro *smiGetFirstMacro(SmiModule *smiModulePtr)
func GetFirstMacro(smiModulePtr *types.SmiModule) *types.SmiMacro {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	macroPtr := modulePtr.Macros.First
	if macroPtr == nil {
		return nil
	}
	return &macroPtr.SmiMacro
}

// SmiMacro *smiGetNextMacro(SmiMacro *smiMacroPtr)
func GetNextMacro(smiMacroPtr *types.SmiMacro) *types.SmiMacro {
	if smiMacroPtr == nil {
		return nil
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	if macroPtr.Next == nil {
		return nil
	}
	return &macroPtr.Next.SmiMacro
}

// SmiModule *smiGetMacroModule(SmiMacro *smiMacroPtr)
func GetMacroModule(smiMacroPtr *types.SmiMacro) *types.SmiModule {
	if smiMacroPtr == nil {
		return nil
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	if macroPtr.Module == nil {
		return nil
	}
	return &macroPtr.Module.SmiModule
}

// int smiGetMacroLine(SmiMacro *smiMacroPtr)
func GetMacroLine(smiMacroPtr *types.SmiMacro) int {
	if smiMacroPtr == nil {
		return 0
	}
	macroPtr := (*internal.Macro)(unsafe.Pointer(smiMacroPtr))
	return macroPtr.Line
}
