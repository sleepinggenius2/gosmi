package models

import (
	"fmt"
)

func GetEnumFormatted(value interface{}, flags Format, enum *Enum) (v Value) {
	intVal, err := ToInt64(value)
	v.Format = flags
	v.Raw = intVal
	if err != nil {
		return
	}
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
