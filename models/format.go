package models

//go:generate enumer -type=Format -autotrimprefix -json

import (
	"fmt"
	"time"

	"github.com/sleepinggenius2/gosmi/types"
)

type Format byte

const (
	FormatNone     Format = 0
	FormatEnumName Format = 1 << iota
	FormatEnumValue
	FormatBits
	FormatString
	FormatUnits
	FormatDurationShort
	FormatAll Format = 0xff & ^FormatUnits
)

func ResolveFormat(formats []Format, defaultFormat ...Format) (format Format) {
	if len(formats) == 0 {
		if len(defaultFormat) == 0 {
			return FormatAll
		}
		return defaultFormat[0]
	}
	for _, f := range formats {
		format |= f
	}
	return
}

type Value struct {
	Format    Format
	Formatted string
	Raw       interface{}
}

func (v Value) Bytes() []byte {
	if b, ok := v.Raw.([]byte); ok {
		return b
	}
	return []byte{}
}

func (v Value) Duration() time.Duration {
	if d, ok := v.Raw.(time.Duration); ok {
		return d
	}
	return time.Duration(0)
}

func (v Value) Int64() int64 {
	if i, ok := v.Raw.(int64); ok {
		return i
	}
	return 0
}

func (v Value) String() string {
	if v.Format == FormatNone {
		return fmt.Sprintf("%v", v.Raw)
	}
	return v.Formatted
}

type ValueFormatter func(interface{}) Value

func (n Node) FormatValue(value interface{}, flags ...Format) Value {
	return n.Type.FormatValue(value, flags...)
}

func (n Node) GetValueFormatter(flags ...Format) ValueFormatter {
	return n.Type.GetValueFormatter(flags...)
}

func (n ScalarNode) FormatValue(value interface{}, flags ...Format) Value {
	return n.Type.FormatValue(value, flags...)
}

func (n ScalarNode) GetValueFormatter(flags ...Format) ValueFormatter {
	return n.Type.GetValueFormatter(flags...)
}

func (n ColumnNode) FormatValue(value interface{}, flags ...Format) Value {
	return n.Type.FormatValue(value, flags...)
}

func (n ColumnNode) GetValueFormatter(flags ...Format) ValueFormatter {
	return n.Type.GetValueFormatter(flags...)
}

func (t Type) FormatValue(value interface{}, flags ...Format) Value {
	formatFlags := ResolveFormat(flags)
	switch t.BaseType {
	case types.BaseTypeOctetString:
		switch t.Name {
		case "IpAddress", "InetAddress", "IpV4orV6Addr":
			return GetInetAddressFormatted(value, formatFlags)
		}
		return GetOctetStringFormatted(value, formatFlags, t.Format)
	case types.BaseTypeBits:
		if t.Enum == nil {
			return GetBitsFormatted(value, formatFlags)
		}
		return GetEnumBitsFormatted(value, formatFlags, t.Enum)
	case types.BaseTypeEnum:
		return GetEnumFormatted(value, formatFlags, t.Enum)
	}
	switch t.Name {
	case "TimeTicks", "TimeInterval", "TimeStamp":
		return GetDurationFormatted(value, formatFlags)
	}
	return GetIntFormatted(value, formatFlags, t.Format)
}

func (t Type) GetValueFormatter(flags ...Format) ValueFormatter {
	formatFlags := ResolveFormat(flags)
	switch t.BaseType {
	case types.BaseTypeOctetString:
		switch t.Name {
		case "IpAddress", "InetAddress", "IpV4orV6Addr":
			return GetInetAddressFormatter(formatFlags)
		}
		return GetOctetStringFormatter(formatFlags, t.Format)
	case types.BaseTypeBits:
		if t.Enum == nil {
			return GetBitsFormatter(formatFlags)
		}
		return GetEnumBitsFormatter(formatFlags, t.Enum)
	case types.BaseTypeEnum:
		return GetEnumFormatter(formatFlags, t.Enum)
	}
	switch t.Name {
	case "TimeTicks", "TimeInterval", "TimeStamp":
		return GetDurationFormatter(formatFlags)
	}
	return GetIntFormatter(formatFlags, t.Format)
}
