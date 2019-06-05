package internal

import (
	"github.com/sleepinggenius2/gosmi/parser"
	"github.com/sleepinggenius2/gosmi/types"
)

const WellKnownModuleName types.SmiIdentifier = "<well-known>"

var (
	WellKnownIdCcitt         types.SmiSubId = 0
	WellKnownIdIso           types.SmiSubId = 1
	WellKnownIdJointIsoCcitt types.SmiSubId = 2
)

type Handle struct {
	Name                 string
	Prev                 *Handle
	Next                 *Handle
	Modules              ModuleMap
	RootNode             *Node
	TypeBits             *Type
	TypeEnum             *Type
	TypeInteger32        *Type
	TypeInteger64        *Type
	TypeObjectIdentifier *Type
	TypeOctetString      *Type
	TypeUnsigned32       *Type
	TypeUnsigned64       *Type
	Flags                Flags
	Paths                []string
	Cache                string
	CacheProg            string
	ErrorLevel           int
	ErrorHandler         types.SmiErrorHandler
}

var smiHandle, firstHandlePtr, lastHandlePtr *Handle

func addHandle(name string) *Handle {
	handlePtr := &Handle{
		Name: name,
		Prev: lastHandlePtr,
	}
	if lastHandlePtr == nil {
		firstHandlePtr = handlePtr
	} else {
		lastHandlePtr.Next = handlePtr
	}
	lastHandlePtr = handlePtr
	return handlePtr
}

func removeHandle(handlePtr *Handle) {
	if handlePtr.Prev != nil {
		handlePtr.Prev.Next = handlePtr.Next
	} else {
		firstHandlePtr = handlePtr.Next
	}
	if handlePtr.Next != nil {
		handlePtr.Next.Prev = handlePtr.Prev
	} else {
		lastHandlePtr = handlePtr.Prev
	}
}

func findHandleByName(name string) *Handle {
	for handlePtr := firstHandlePtr; handlePtr != nil; handlePtr = handlePtr.Next {
		if handlePtr.Name == name {
			return handlePtr
		}
	}
	return nil
}

func SetErrorHandler(smiErrorHandler types.SmiErrorHandler) {
	smiHandle.ErrorHandler = smiErrorHandler
}

func SetSeverity(pattern string, severity int) {}

func SetErrorLevel(level int) {}

func GetFlags() int { return 0 }

func SetFlags(userflags int) {}

func Initialized() bool {
	return smiHandle != nil
}

func GetFirstModule() *Module {
	if smiHandle == nil {
		return nil
	}
	return smiHandle.Modules.First
}

func Root() *Node {
	if smiHandle == nil {
		return nil
	}
	return smiHandle.RootNode
}

func oidFromSubId(subId types.SmiSubId) parser.Oid {
	return parser.Oid{
		SubIdentifiers: []parser.SubIdentifier{
			parser.SubIdentifier{Number: &subId},
		},
	}
}

func createBaseType(module *Module, baseType types.BaseType) *Type {
	return &Type{
		SmiType: types.SmiType{
			Name:     types.SmiIdentifier(baseType.String()),
			BaseType: baseType,
			Decl:     types.DeclImplicitType,
		},
		Module: module,
	}
}

func initData() bool {
	smiHandle.RootNode = &Node{Flags: FlagRoot, Oid: types.Oid{}}

	wellKnownModule := &Module{
		SmiModule: types.SmiModule{
			Name: WellKnownModuleName,
		},
		PrefixNode: smiHandle.RootNode,
	}

	// Create ccitt well-known node
	ccitt := &Object{
		SmiNode: types.SmiNode{
			Name:     "ccitt",
			Decl:     types.DeclImplObject,
			NodeKind: types.NodeNode,
		},
		Module: wellKnownModule,
	}
	wellKnownModule.Objects.AddWithOid(ccitt, oidFromSubId(WellKnownIdCcitt))

	// Create iso well-known node
	iso := &Object{
		SmiNode: types.SmiNode{
			Name:     "iso",
			Decl:     types.DeclImplObject,
			NodeKind: types.NodeNode,
		},
		Module: wellKnownModule,
	}
	wellKnownModule.Objects.AddWithOid(iso, oidFromSubId(WellKnownIdIso))

	// Create joint-iso-ccitt well-known node
	jointIsoCcitt := &Object{
		SmiNode: types.SmiNode{
			Name:     "joint-iso-ccitt",
			Decl:     types.DeclImplObject,
			NodeKind: types.NodeNode,
		},
		Module: wellKnownModule,
	}
	wellKnownModule.Objects.AddWithOid(jointIsoCcitt, oidFromSubId(WellKnownIdJointIsoCcitt))

	smiHandle.Modules.Add(wellKnownModule)

	smiHandle.TypeBits = createBaseType(wellKnownModule, types.BaseTypeBits)
	smiHandle.TypeEnum = createBaseType(wellKnownModule, types.BaseTypeEnum)
	smiHandle.TypeInteger32 = createBaseType(wellKnownModule, types.BaseTypeInteger32)
	smiHandle.TypeInteger64 = createBaseType(wellKnownModule, types.BaseTypeInteger64)
	smiHandle.TypeObjectIdentifier = createBaseType(wellKnownModule, types.BaseTypeObjectIdentifier)
	smiHandle.TypeOctetString = createBaseType(wellKnownModule, types.BaseTypeOctetString)
	smiHandle.TypeUnsigned32 = createBaseType(wellKnownModule, types.BaseTypeUnsigned32)
	smiHandle.TypeUnsigned64 = createBaseType(wellKnownModule, types.BaseTypeUnsigned64)

	return true
}

func freeData() {
}
