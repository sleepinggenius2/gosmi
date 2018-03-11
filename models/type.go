package models

import (
	"fmt"
	"sort"

	"github.com/sleepinggenius2/gosmi/types"
)

type Enum struct {
	BaseType types.BaseType
	Values   EnumValues
}

type EnumValues map[int64]string

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

// Keys returns a sorted list of all keys
func (values EnumValues) Keys() []int64 {
	keys := make([]int64, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
