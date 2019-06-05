package smi

import (
	"fmt"
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// char *smiLoadModule(const char *module)
func LoadModule(module string) string {
	checkInit()
	modulePtr, err := internal.GetModule(module)
	if err != nil {
		fmt.Println(err)
	}
	if modulePtr == nil {
		return ""
	}
	return modulePtr.Name.String()
}

// int smiIsLoaded(const char *module)
func IsLoaded(module string) bool {
	checkInit()
	return internal.FindModuleByName(module) != nil
}

// SmiModule *smiGetModule(const char *module)
func GetModule(module string) *types.SmiModule {
	if module == "" {
		return nil
	}
	modulePtr, _ := internal.GetModule(module)
	if modulePtr == nil {
		return nil
	}
	return &modulePtr.SmiModule
}

// SmiModule *smiGetFirstModule(void)
func GetFirstModule() *types.SmiModule {
	modulePtr := internal.GetFirstModule()
	if modulePtr == nil {
		return nil
	}
	return &modulePtr.SmiModule
}

// SmiModule *smiGetNextModule(SmiModule *smiModulePtr)
func GetNextModule(smiModulePtr *types.SmiModule) *types.SmiModule {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.Next == nil {
		return nil
	}
	return &modulePtr.Next.SmiModule

}

// SmiNode *smiGetModuleIdentityNode(SmiModule *smiModulePtr)
func GetModuleIdentityNode(smiModulePtr *types.SmiModule) *types.SmiNode {
	if smiModulePtr == nil {
		return nil
	}
	modulePtr := (*internal.Module)(unsafe.Pointer(smiModulePtr))
	if modulePtr.Identity == nil {
		return nil
	}
	return modulePtr.Identity.GetSmiNode()
}
