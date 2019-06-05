package smi

import (
	"math"
	"unsafe"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

// SmiRange *smiGetFirstRange(SmiType *smiTypePtr)
func GetFirstRange(smiTypePtr *types.SmiType) *types.SmiRange {
	if smiTypePtr == nil {
		return nil
	}
	typePtr := (*internal.Type)(unsafe.Pointer(smiTypePtr))
	if typePtr.List == nil || typePtr.List.Ptr == nil || typePtr.BaseType == types.BaseTypeEnum || typePtr.BaseType == types.BaseTypeBits {
		return nil
	}
	return &typePtr.List.Ptr.(*internal.Range).SmiRange
}

// SmiRange *smiGetNextRange(SmiRange *smiRangePtr)
func GetNextRange(smiRangePtr *types.SmiRange) *types.SmiRange {
	if smiRangePtr == nil {
		return nil
	}
	rangePtr := (*internal.Range)(unsafe.Pointer(smiRangePtr))
	if rangePtr.Type == nil || rangePtr.Type.BaseType == types.BaseTypeEnum || rangePtr.Type.BaseType == types.BaseTypeBits {
		return nil
	}
	if rangePtr.List == nil || rangePtr.List.Next == nil || rangePtr.List.Next.Ptr == nil {
		return nil
	}
	return &rangePtr.List.Next.Ptr.(*internal.Range).SmiRange
}

// int smiGetMinMaxRange(SmiType *smiType, SmiValue *min, SmiValue *max)
func GetMinMaxRange(smiType *types.SmiType) *types.SmiRange {
	if smiType == nil {
		return nil
	}
	currRange := GetFirstRange(smiType)
	if currRange == nil {
		return nil
	}
	baseType := currRange.MinValue.BaseType
	rangePtr := new(types.SmiRange)
	rangePtr.MinValue.BaseType = baseType
	rangePtr.MaxValue.BaseType = baseType
	switch baseType {
	case types.BaseTypeInteger32:
		rangePtr.MinValue.Value = int32(math.MaxInt32)
		rangePtr.MaxValue.Value = int32(math.MinInt32)
	case types.BaseTypeInteger64:
		rangePtr.MinValue.Value = int64(math.MaxInt64)
		rangePtr.MaxValue.Value = int64(math.MinInt64)
	case types.BaseTypeUnsigned32:
		rangePtr.MinValue.Value = uint32(math.MaxUint32)
		rangePtr.MaxValue.Value = 0
	case types.BaseTypeUnsigned64:
		rangePtr.MinValue.Value = uint64(math.MaxUint64)
		rangePtr.MaxValue.Value = 0
	default:
		return nil
	}
	for ; currRange != nil; currRange = GetNextRange(currRange) {
		switch baseType {
		case types.BaseTypeInteger32:
			if currRange.MinValue.Value.(int32) < rangePtr.MinValue.Value.(int32) {
				rangePtr.MinValue.Value = currRange.MinValue.Value
			}
			if currRange.MaxValue.Value.(int32) > rangePtr.MaxValue.Value.(int32) {
				rangePtr.MaxValue.Value = currRange.MaxValue.Value
			}
		case types.BaseTypeInteger64:
			if currRange.MinValue.Value.(int64) < rangePtr.MinValue.Value.(int64) {
				rangePtr.MinValue.Value = currRange.MinValue.Value
			}
			if currRange.MaxValue.Value.(int64) > rangePtr.MaxValue.Value.(int64) {
				rangePtr.MaxValue.Value = currRange.MaxValue.Value
			}
		case types.BaseTypeUnsigned32:
			if currRange.MinValue.Value.(int32) < rangePtr.MinValue.Value.(int32) {
				rangePtr.MinValue.Value = currRange.MinValue.Value
			}
			if currRange.MaxValue.Value.(int32) > rangePtr.MaxValue.Value.(int32) {
				rangePtr.MaxValue.Value = currRange.MaxValue.Value
			}
		case types.BaseTypeUnsigned64:
			if currRange.MinValue.Value.(uint64) < rangePtr.MinValue.Value.(uint64) {
				rangePtr.MinValue.Value = currRange.MinValue.Value
			}
			if currRange.MaxValue.Value.(uint64) > rangePtr.MaxValue.Value.(uint64) {
				rangePtr.MaxValue.Value = currRange.MaxValue.Value
			}
		}
	}
	return rangePtr
}
