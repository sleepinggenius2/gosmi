package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"encoding/binary"
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
	smiType     *C.struct_SmiType
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

func (t *Type) getEnum() {
	if t.BaseType == types.BaseTypeUnknown || !(t.BaseType == types.BaseTypeEnum || t.BaseType == types.BaseTypeBits) {
		return
	}

	smiNamedNumber := C.smiGetFirstNamedNumber(t.smiType)
	if smiNamedNumber == nil {
		return
	}

	enum := Enum{
		BaseType: types.BaseType(smiNamedNumber.value.basetype),
	}
	for ; smiNamedNumber != nil; smiNamedNumber = C.smiGetNextNamedNumber(smiNamedNumber) {
		namedNumber := NamedNumber{
			Name:  C.GoString(smiNamedNumber.name),
			Value: convertValue(smiNamedNumber.value),
		}
		enum.Values = append(enum.Values, namedNumber)
	}
	t.Enum = &enum
	return
}

func (t Type) GetModule() (module Module) {
	smiModule := C.smiGetTypeModule(t.smiType)
	return CreateModule(smiModule)
}

func (t *Type) getRanges() {
	if t.BaseType == types.BaseTypeUnknown {
		return
	}

	ranges := make([]Range, 0)
	for smiRange := C.smiGetFirstRange(t.smiType); smiRange != nil; smiRange = C.smiGetNextRange(smiRange) {
		r := Range{
			BaseType: types.BaseType(smiRange.minValue.basetype),
			MinValue: convertValue(smiRange.minValue),
			MaxValue: convertValue(smiRange.maxValue),
		}
		ranges = append(ranges, r)
	}
	t.Ranges = ranges
}

func (t Type) String() string {
	typeStr := t.Name
	if t.BaseType.String() != typeStr {
		typeStr += "<" + t.BaseType.String() + ">"
	}
	return fmt.Sprintf("Type[%s Status=%s, Format=%s, Units=%s]", typeStr, t.Status, t.Format, t.Units)
}

func (t Type) GetRaw() (outType *C.struct_SmiType) {
	return t.smiType
}

func (t *Type) SetRaw(smiType *C.struct_SmiType) {
	t.smiType = smiType
}

func CreateType(smiType *C.struct_SmiType) (outType Type) {
	if smiType == nil {
		return
	}

	outType.SetRaw(smiType)
	outType.BaseType = types.BaseType(smiType.basetype)

	if smiType.name == nil {
		smiType = C.smiGetParentType(smiType)
	}

	outType.Decl = types.Decl(smiType.decl)
	outType.Description = C.GoString(smiType.description)
	outType.Format = C.GoString(smiType.format)
	outType.Name = C.GoString(smiType.name)
	outType.Reference = C.GoString(smiType.reference)
	outType.Status = types.Status(smiType.status)
	outType.Units = C.GoString(smiType.units)

	outType.getEnum()
	outType.getRanges()

	return
}

func CreateTypeFromNode(smiNode *C.struct_SmiNode) (outType *Type) {
	smiType := C.smiGetNodeType(smiNode)

	if smiType == nil {
		return
	}

	tempType := CreateType(smiType)
	outType = &tempType

	if format := C.GoString(smiNode.format); format != "" {
		outType.Format = format
	}
	if units := C.GoString(smiNode.units); units != "" {
		outType.Units = units
	}

	return
}

func convertValue(value C.struct_SmiValue) (outValue int64) {
	switch types.BaseType(value.basetype) {
	case types.BaseTypeInteger32:
		outValue = int64(int32(binary.LittleEndian.Uint32(value.value[:4])))
	case types.BaseTypeInteger64:
		outValue = int64(binary.LittleEndian.Uint64(value.value[:8]))
	case types.BaseTypeUnsigned32:
		outValue = int64(binary.LittleEndian.Uint32(value.value[:4]))
	case types.BaseTypeUnsigned64:
		outValue = int64(binary.LittleEndian.Uint64(value.value[:8]))
	}
	return
}
