package models

import (
	"github.com/min-oc/gosmi/types"
)

type Node struct {
	Access      types.Access
	Decl        types.Decl
	Description string
	Kind        types.NodeKind
	Name        string
	Oid         types.Oid
	OidLen      int
	Status      types.Status
	Type        *Type
}
