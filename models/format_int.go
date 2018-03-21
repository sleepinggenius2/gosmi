package models

import (
	"fmt"
	"strconv"
)

func GetIntFormatted(value interface{}, flags Format, format string) Value {
	intVal := ToInt64(value)
	var formatted string
	if flags != FormatNone {
		formatted = IntegerDisplayHint(format, intVal)
	}
	return Value{
		Format:    flags,
		Formatted: formatted,
		Raw:       intVal,
	}
}

func GetIntFormatter(flags Format, format string) (f ValueFormatter) {
	return func(value interface{}) Value {
		return GetIntFormatted(value, flags, format)
	}
}

func IntegerDisplayHint(format string, value int64) (formatted string) {
	if len(format) == 0 {
		return fmt.Sprintf("%d", value)
	}
	switch format[0] {
	case 'b':
		formatted = fmt.Sprintf("%b", value)
	case 'd':
		formatted = fmt.Sprintf("%d", value)
		if len(format) < 3 {
			break
		}
		decimals, err := strconv.Atoi(format[2:])
		if err != nil || decimals < 1 {
			break
		}
		formattedLen := len(formatted)
		offset := 0
		if formatted[0] == '-' {
			offset = 1
		}
		if formattedLen-offset <= decimals {
			formatStr := "0.%0" + format[2:] + "s"
			formatted = formatted[:offset] + fmt.Sprintf(formatStr, formatted[offset:])
			break
		}
		formatted = formatted[:formattedLen-decimals] + "." + formatted[formattedLen-decimals:]
	case 'o':
		formatted = fmt.Sprintf("%o", value)
	case 'x':
		formatted = fmt.Sprintf("%x", value)
	default:
		formatted = fmt.Sprintf("%d", value)
	}
	return
}
