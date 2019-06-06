package internal

import (
	"sort"
	"strconv"

	"github.com/sleepinggenius2/gosmi/parser"
	"github.com/sleepinggenius2/gosmi/types"
)

type Type struct {
	types.SmiType
	Module *Module
	Parent *Type
	List   *List
	Flags  Flags
	Prev   *Type
	Next   *Type
	Line   int

	lastList *List
}

func (x *Type) AddList(list *List) {
	if x.lastList == nil {
		x.List = list
	} else {
		x.lastList.Next = list
	}
	x.lastList = list
}

func (x *Type) AddRange(min types.SmiValue, max types.SmiValue) {
	list := &List{}
	list.Ptr = &Range{
		SmiRange: types.SmiRange{min, max},
		Type:     x,
		List:     list,
	}
	x.AddList(list)
}

func (x *Type) AddNamedNumber(name types.SmiIdentifier, value types.SmiValue) {
	list := &List{}
	list.Ptr = &NamedNumber{
		SmiNamedNumber: types.SmiNamedNumber{name, value},
		Type:           x,
		List:           list,
	}
	x.AddList(list)
}

type TypeMap struct {
	First *Type

	last *Type
	m    map[types.SmiIdentifier]*Type
}

func (x *TypeMap) Add(t *Type) {
	t.Prev = x.last
	if x.First == nil {
		x.First = t
	} else {
		x.last.Next = t
	}
	x.last = t

	if x.m == nil {
		x.m = make(map[types.SmiIdentifier]*Type)
	}
	x.m[t.Name] = t
}

func (x *TypeMap) Get(name types.SmiIdentifier) *Type {
	if x.m == nil {
		return nil
	}
	return x.m[name]
}

func (x *TypeMap) GetName(name string) *Type {
	return x.Get(types.SmiIdentifier(name))
}

type NamedNumber struct {
	types.SmiNamedNumber
	Type *Type
	List *List
}

type Range struct {
	types.SmiRange
	Type *Type
	List *List
}

func GetBaseTypeFromSyntax(syntax parser.SyntaxType) *Type {
	switch syntax.Name {
	case "BITS":
		return smiHandle.TypeBits
	case "INTEGER":
		if syntax.SubType == nil || len(syntax.SubType.Integer) == 0 {
			return smiHandle.TypeInteger32
		}

		// Assuming the ranges are in order
		minValue := syntax.SubType.Integer[0].Start
		maxValue := syntax.SubType.Integer[len(syntax.SubType.Integer)-1].End
		if maxValue == "" {
			maxValue = syntax.SubType.Integer[len(syntax.SubType.Integer)-1].Start
		}

		// The parser should guarantee that there is at least 1 digit
		if minValue[0] == '-' {
			if len(minValue) > 11 || minValue[1:] > "2147483648" {
				return smiHandle.TypeInteger64
			}
			return smiHandle.TypeInteger32
		} else {
			maxLen := len(maxValue)
			// Check for BinString or HexString
			if maxValue[0] == '\'' {
				if maxValue[maxLen-1] == 'H' && maxLen > 11 { // 8 hex digits + 3 wrapper chars
					return smiHandle.TypeUnsigned64
				} else if maxValue[maxLen-1] == 'B' && maxLen > 35 { // 32 binary digits + 3 wrapper chars
					return smiHandle.TypeUnsigned64
				}
			} else if maxLen > 10 || maxValue > "4294967295" {
				return smiHandle.TypeUnsigned64
			}
			return smiHandle.TypeUnsigned32
		}
	case "OBJECT IDENTIFIER":
		return smiHandle.TypeObjectIdentifier
	case "OCTET STRING":
		return smiHandle.TypeOctetString
	}
	return nil
}

func getValueInt(value string, bits int) int64 {
	if value == "" {
		return 0
	}
	if value[0] != '\'' {
		i, _ := strconv.ParseInt(value, 10, bits)
		return i
	}
	if len(value) < 4 {
		return 0
	}
	var base int
	switch value[len(value)-1] {
	case 'B':
		base = 2
	case 'H':
		base = 16
	default:
		return 0
	}
	i, _ := strconv.ParseInt(value[1:len(value)-2], base, bits)
	return i
}

func getValueUint(value string, bits int) uint64 {
	if value == "" {
		return 0
	}
	if value[0] != '\'' {
		i, _ := strconv.ParseUint(value, 10, bits)
		return i
	}
	if len(value) < 4 {
		return 0
	}
	var base int
	switch value[len(value)-1] {
	case 'B':
		base = 2
	case 'H':
		base = 16
	default:
		return 0
	}
	i, _ := strconv.ParseUint(value[1:len(value)-2], base, bits)
	return i
}

func GetValue(value string, baseType types.BaseType) types.SmiValue {
	v := types.SmiValue{BaseType: baseType}
	switch baseType {
	case types.BaseTypeInteger32:
		v.Value = GetValueInt32(value)
	case types.BaseTypeInteger64:
		v.Value = GetValueInt64(value)
	case types.BaseTypeUnsigned32:
		v.Value = GetValueUint32(value)
	case types.BaseTypeUnsigned64:
		v.Value = GetValueUint64(value)
	}
	return v
}

func GetValueInt32(value string) int32   { return int32(getValueInt(value, 32)) }
func GetValueUint32(value string) uint32 { return uint32(getValueUint(value, 32)) }
func GetValueInt64(value string) int64   { return getValueInt(value, 64) }
func GetValueUint64(value string) uint64 { return getValueUint(value, 64) }

func intStringLess(a, b string) bool {
	if a[0] == '-' {
		if b[0] == '-' {
			if len(a) < len(b) {
				return false
			}
			return len(a) > len(b) || a[1:] > b[1:]
		}
		return true
	} else if b[0] == '-' || len(a) > len(b) {
		return false
	}
	return len(a) < len(b) || a < b
}

func namedNumberSort(namedNumbers []parser.NamedNumber) {
	sort.Slice(namedNumbers, func(i, j int) bool {
		return intStringLess(namedNumbers[i].Value, namedNumbers[j].Value)
	})
}

func rangeSort(ranges []parser.Range) {
	sort.Slice(ranges, func(i, j int) bool {
		return intStringLess(ranges[i].Start, ranges[j].Start)
	})
}
