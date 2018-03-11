package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import (
	"encoding/binary"

	"github.com/sleepinggenius2/gosmi/models"
	"github.com/sleepinggenius2/gosmi/types"
)

type SmiType struct {
	models.Type
	smiType *C.struct_SmiType
}

func (t *SmiType) getEnum() {
	if t.BaseType == types.BaseTypeUnknown || !(t.BaseType == types.BaseTypeEnum || t.BaseType == types.BaseTypeBits) {
		return
	}

	smiNamedNumber := C.smiGetFirstNamedNumber(t.smiType)
	if smiNamedNumber == nil {
		return
	}

	t.Enum = &models.Enum{
		BaseType: types.BaseType(smiNamedNumber.value.basetype),
		Values:   make(models.EnumValues),
	}
	for ; smiNamedNumber != nil; smiNamedNumber = C.smiGetNextNamedNumber(smiNamedNumber) {
		t.Enum.Values[convertValue(smiNamedNumber.value)] = C.GoString(smiNamedNumber.name)
	}
	return
}

func (t SmiType) GetModule() (module SmiModule) {
	smiModule := C.smiGetTypeModule(t.smiType)
	return CreateModule(smiModule)
}

func (t *SmiType) getRanges() {
	if t.BaseType == types.BaseTypeUnknown {
		return
	}

	ranges := make([]models.Range, 0)
	for smiRange := C.smiGetFirstRange(t.smiType); smiRange != nil; smiRange = C.smiGetNextRange(smiRange) {
		r := models.Range{
			BaseType: types.BaseType(smiRange.minValue.basetype),
			MinValue: convertValue(smiRange.minValue),
			MaxValue: convertValue(smiRange.maxValue),
		}
		ranges = append(ranges, r)
	}
	t.Ranges = ranges
}

func (t SmiType) String() string {
	return t.Type.String()
}

func (t SmiType) GetRaw() (outType *C.struct_SmiType) {
	return t.smiType
}

func (t *SmiType) SetRaw(smiType *C.struct_SmiType) {
	t.smiType = smiType
}

func CreateType(smiType *C.struct_SmiType) (outType SmiType) {
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

func CreateTypeFromNode(smiNode *C.struct_SmiNode) (outType *SmiType) {
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
