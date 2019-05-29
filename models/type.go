package models

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/sleepinggenius2/gosmi/types"
)

type Enum struct {
	BaseType types.BaseType
	Values   []NamedNumber
	valueMap map[int64]string
	rw       sync.RWMutex
}

func (e *Enum) initValueMap() {
	e.rw.RLock()
	if e.valueMap != nil {
		e.rw.RUnlock()
		return
	}
	e.rw.RUnlock()
	e.rw.Lock()
	e.valueMap = make(map[int64]string, len(e.Values))
	for _, value := range e.Values {
		e.valueMap[value.Value] = value.Name
	}
	e.rw.Unlock()
}

func (e *Enum) Name(value int64) string {
	e.initValueMap()
	e.rw.RLock()
	name, ok := e.valueMap[value]
	e.rw.RUnlock()
	if !ok {
		return "unknown"
	}
	return name
}

func (e *Enum) Value(name string) (int64, error) {
	e.initValueMap()
	e.rw.RLock()
	defer e.rw.RUnlock()
	for k, v := range e.valueMap {
		if v == name {
			return k, nil
		}
	}
	return 0, errors.New("Unknown enum name: " + name)
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

func (t Type) indexValueEnum(value interface{}) ([]uint32, error) {
	var intVal int64
	var err error
	if strVal, ok := value.(string); ok && t.Enum != nil {
		intVal, err = t.Enum.Value(strVal)
	} else {
		intVal, err = ToInt64(value)
	}
	if err != nil {
		return nil, err
	}
	if intVal < 0 || intVal > 0xffffffff {
		return nil, errors.New("Integer value outside of range")
	}
	return []uint32{uint32(intVal)}, nil
}

func (t Type) indexValueInteger(value interface{}) ([]uint32, error) {
	intVal, err := ToInt64(value)
	if err != nil {
		return nil, err
	}
	if intVal < 0 || intVal > 0xffffffff {
		return nil, errors.New("Integer value outside of range")
	}
	return []uint32{uint32(intVal)}, nil
}

func (t Type) indexValueObjectIdentifier(value []uint32, implied bool) ([]uint32, error) {
	var offset int
	if !implied {
		offset = 1
	}
	ret := make([]uint32, len(value)+offset)
	if !implied {
		ret[0] = uint32(len(value))
	}
	copy(ret[offset:], value)
	return ret, nil
}

func (t Type) indexValueOctetString(value interface{}, implied bool) ([]uint32, error) {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return nil, errors.New("Invalid octet string value")
	}
	var ret []uint32
	var offset int
	if implied {
		ret = make([]uint32, len(bytes))
	} else {
		ret = make([]uint32, len(bytes)+1)
		ret[0] = uint32(len(bytes))
		offset = 1
	}
	for i, b := range bytes {
		ret[i+offset] = uint32(b)
	}
	return ret, nil
}

func (t Type) IndexValue(value interface{}, implied bool) ([]uint32, error) {
	switch t.BaseType {
	case types.BaseTypeEnum:
		return t.indexValueEnum(value)
	case types.BaseTypeInteger32, types.BaseTypeUnsigned32:
		return t.indexValueInteger(value)
	case types.BaseTypeObjectIdentifier:
		switch v := value.(type) {
		case []uint32:
			return t.indexValueObjectIdentifier(v, implied)
		case Oid:
			return t.indexValueObjectIdentifier(v, implied)
		case string:
			oid, err := OidFromString(v)
			if err != nil {
				return nil, err
			}
			return t.indexValueObjectIdentifier(oid, implied)
		}
		return nil, errors.New("Invalid object identifier value")
	case types.BaseTypeOctetString:
		return t.indexValueOctetString(value, implied)
	}
	return nil, errors.Errorf("Invalid base type: %v", t.BaseType)
}
