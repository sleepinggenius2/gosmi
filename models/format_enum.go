package models

import (
	"fmt"

	"github.com/sleepinggenius2/gosnmp"
)

func GetEnumFormatted(value interface{}, flags Format, enum *Enum) (v Value) {
	v.Format = flags
	v.Raw = value
	if flags == FormatNone {
		return
	}
	var intVal int64
	switch tempVal := value.(type) {
	case int64:
		intVal = tempVal
	default:
		intVal = gosnmp.ToBigInt(tempVal).Int64()
	}
	v.Raw = intVal
	if flags&FormatEnumName != 0 {
		v.Formatted = enum.Name(intVal)
		if flags&FormatEnumValue != 0 {
			v.Formatted += fmt.Sprintf("(%d)", intVal)
		}
	} else if flags&FormatEnumValue != 0 {
		v.Formatted = fmt.Sprintf("%d", intVal)
	}
	return
}

func GetEnumFormatter(flags Format, enum *Enum) (f ValueFormatter) {
	return func(value interface{}) Value {
		return GetEnumFormatted(value, flags, enum)
	}
}
