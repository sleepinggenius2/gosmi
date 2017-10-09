package models

import (
	"fmt"

	"github.com/sleepinggenius2/gosmi/types"
)

type Enum struct {
	BaseType types.BaseType
	Values   []NamedNumber
}

type NamedNumber struct {
	Name  string
	Value int64
}

type Range struct {
	BaseType types.BaseType
	MinValue int64
	MaxValue int64
}

type Type struct {
	BaseType    types.BaseType
	Decl        types.Decl
	Description string
	Enum        *Enum
	Format      string
	Name        string
	Ranges      []Range
	Reference   string
	Status      types.Status
	Units       string
}

func (t Type) String() string {
	typeStr := t.Name
	if t.BaseType.String() != typeStr {
		typeStr += "<" + t.BaseType.String() + ">"
	}
	return fmt.Sprintf("Type[%s Status=%s, Format=%s, Units=%s]", typeStr, t.Status, t.Format, t.Units)
}
