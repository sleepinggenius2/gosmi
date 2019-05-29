package models

import (
	"fmt"
	"strconv"
	"strings"
)

func GetIntFormatted(value interface{}, flags Format, format string) Value {
	var formatted string
	intVal, err := ToInt64(value)
	if err == nil && flags != FormatNone {
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
		return strconv.FormatInt(value, 10)
	}
	switch format[0] {
	case 'b':
		formatted = fmt.Sprintf("%b", value)
	case 'd':
		formatted = strconv.FormatInt(value, 10)
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
			zeros := decimals - formattedLen + offset
			formatted = formatted[:offset] + "0." + strings.Repeat("0", zeros) + formatted[offset:]
		} else {
			formatted = formatted[:formattedLen-decimals] + "." + formatted[formattedLen-decimals:]
		}
	case 'o':
		formatted = strconv.FormatInt(value, 8)
	case 'x':
		formatted = strconv.FormatInt(value, 16)
	default:
		formatted = strconv.FormatInt(value, 10)
	}
	return
}
