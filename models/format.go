package models

//go:generate enumer -type=Format -autotrimprefix -json

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sleepinggenius2/gosmi/types"
	"github.com/sleepinggenius2/gosnmp"
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

type Value struct {
	Format    Format
	Formatted string
	Raw       interface{}
}

func (v Value) String() string {
	if v.Format == FormatNone {
		return fmt.Sprintf("%v", v.Raw)
	}
	return v.Formatted
}

func (v Value) Int64() int64 {
	if i, ok := v.Raw.(int64); ok {
		return i
	}
	return 0
}

type ValueFormatter func(interface{}) Value

func (n Node) GetValueFormatter(flags ...Format) (f ValueFormatter) {
	return n.Type.GetValueFormatter(flags...)
}

func (t Type) GetValueFormatter(flags ...Format) (f ValueFormatter) {
	formatFlags := FormatNone
	for _, flag := range flags {
		formatFlags |= flag
	}
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
		return GetEnumBitsFormatter(formatFlags, t.Enum.Values)
	case types.BaseTypeEnum:
		return GetEnumFormatter(formatFlags, t.Enum.Values)
	}
	switch t.Name {
	case "TimeTicks", "TimeInterval", "TimeStamp":
		return GetDurationFormatter(formatFlags)
	}
	return GetIntFormatter(formatFlags, t.Format)
}

func GetBitsFormatter(flags Format) (f ValueFormatter) {
	return func(value interface{}) (v Value) {
		v.Format = flags
		v.Raw = value
		if flags&FormatBits != 0 {
			if bytes, ok := value.([]byte); ok {
				v.Formatted = fmt.Sprintf("% X", bytes)
			}
		}
		return
	}
}

func GetDurationFormatter(flags Format) (f ValueFormatter) {
	return func(value interface{}) (v Value) {
		duration := time.Duration(int64(value.(int)) * 1e7)
		v.Format = flags
		v.Raw = duration
		if flags&FormatDurationShort > 0 {
			v.Formatted = DurationFormat(duration)
		} else {
			v.Formatted = DurationFormatLong(duration)
		}
		return
	}
}

func GetEnumFormatter(flags Format, values []NamedNumber) (f ValueFormatter) {
	enums := make(map[int64]string, len(values))
	for _, enum := range values {
		enums[enum.Value] = enum.Name
	}
	return func(value interface{}) (v Value) {
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
			enumName, ok := enums[intVal]
			if !ok {
				enumName = "unknown"
			}
			v.Formatted = enumName
			if flags&FormatEnumValue != 0 {
				v.Formatted += fmt.Sprintf("(%d)", intVal)
			}
		} else if flags&FormatEnumValue != 0 {
			v.Formatted = fmt.Sprintf("%d", intVal)
		}
		return
	}
}

func GetEnumBitsFormatter(flags Format, values []NamedNumber) (f ValueFormatter) {
	enums := make(map[uint64]string, len(values))
	for _, enum := range values {
		enums[uint64(enum.Value)] = enum.Name
	}
	return func(value interface{}) (v Value) {
		v.Format = flags
		v.Raw = value
		if flags == FormatNone {
			return
		}
		octets := value.([]byte)
		if flags&FormatBits != 0 {
			v.Formatted = fmt.Sprintf("% X", octets)
		}
		if (flags&FormatEnumName)+(flags&FormatEnumValue) == 0 {
			return
		}
		bitsFormatted := make([]string, 0, 8*len(octets))
		for i, octet := range octets {
			for j := 7; j >= 0; j-- {
				if octet&(1<<uint(j)) != 0 {
					bit := uint64(8*i + (7 - j))
					var bitFormatted string
					if flags&FormatEnumName != 0 {
						enumName, ok := enums[bit]
						if !ok {
							enumName = "unknown"
						}
						bitFormatted = enumName
						if flags&FormatEnumValue != 0 || !ok {
							bitFormatted += "(" + fmt.Sprintf("%d", bit) + ")"
						}
					} else if flags&FormatEnumValue != 0 {
						bitFormatted = fmt.Sprintf("%d", bit)
					}
					bitsFormatted = append(bitsFormatted, bitFormatted)
				}
			}
		}
		if v.Formatted == "" {
			v.Formatted = strings.Join(bitsFormatted, " ")
		} else {
			v.Formatted += "[" + strings.Join(bitsFormatted, " ") + "]"
		}
		return
	}
}

func getBytes(value interface{}) (bytes []byte, ok bool) {
	switch val := value.(type) {
	case []int:
		bytes = make([]byte, len(val))
		for i, b := range val {
			bytes[i] = byte(b)
		}
		ok = true
	case []byte:
		bytes = val
		ok = true
	}
	return
}

func GetInetAddressFormatter(flags Format) (f ValueFormatter) {
	return func(value interface{}) (v Value) {
		v.Format = flags
		bytes, ok := getBytes(value)
		if !ok {
			v.Raw = value
			return
		}
		v.Raw = bytes
		if flags&FormatString == 0 {
			return
		}
		numBytes := len(bytes)
		if numBytes == 5 || numBytes == 17 {
			bytes = bytes[1:]
			numBytes--
		}
		var format string
		if numBytes == 4 {
			format = "1d.1d.1d.1d"
		} else {
			format = "2x:2x:2x:2x:2x:2x:2x:2x%4d"
		}
		v.Formatted = StringDisplayHint(format, bytes)
		return
	}
}

func GetIntFormatter(flags Format, format string) (f ValueFormatter) {
	return func(value interface{}) (v Value) {
		var intVal int64
		switch tempVal := value.(type) {
		case int64:
			intVal = tempVal
		default:
			intVal = gosnmp.ToBigInt(tempVal).Int64()
		}
		return Value{
			Format:    flags,
			Formatted: IntegerDisplayHint(format, intVal),
			Raw:       intVal,
		}
	}
}

func GetOctetStringFormatter(flags Format, format string) (f ValueFormatter) {
	return func(value interface{}) (v Value) {
		v.Format = flags
		var bytes []byte
		switch val := value.(type) {
		case []int:
			bytes = make([]byte, len(val))
			for i, b := range val {
				bytes[i] = byte(b)
			}
		case []byte:
			bytes = val
		default:
			v.Raw = val
			if flags&FormatString != 0 {
				v.Formatted = fmt.Sprintf("%v", val)
			}
			return
		}
		v.Raw = bytes
		if flags&FormatString == 0 {
			return
		}
		switch format {
		case "InetAddress":
			if len(bytes) > 0 && bytes[0] == 0x04 {
				format = "*1d."
			} else {
				format = "*1x:"
			}
		case "IpV4orV6Addr":
			if len(bytes) == 4 {
				format = "1d."
			} else {
				format = "1x:"
			}
		}
		v.Formatted = StringDisplayHint(format, bytes)
		return
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

func StringDisplayHint(format string, value []byte) (formatted string) {
	if len(format) < 2 {
		return fmt.Sprintf("% X", value)
	}
	formats := parseHint(format)
	formatsLen := len(formats) - 1
	//fmt.Printf("%s = %+v\n", format, formats)
	var (
		i, repeatCount, lengthCount int
		currValue                   uint64
	)
	currFormat := formats[i]
	for j, c := range value {
		if currFormat.Repeat && repeatCount == 0 {
			repeatCount = int(c)
			continue
		}
		if lengthCount == 0 {
			lengthCount = currFormat.Length
			currValue = 0
		}
		lengthCount--
		if currFormat.Numeric {
			currValue = currValue<<8 + uint64(c)
		} else {
			formatted += fmt.Sprintf(currFormat.Format, c)
		}
		if lengthCount == 0 {
			if currFormat.Numeric {
				formatted += fmt.Sprintf(currFormat.Format, currValue)
			}
			if currFormat.Separator > 0 && j != len(value)-1 {
				formatted += string([]byte{currFormat.Separator})
			}
			if currFormat.Repeat {
				repeatCount--
				if repeatCount == 0 {
					if currFormat.Terminator > 0 && j != len(value)-1 {
						formatted += string([]byte{currFormat.Terminator})
					}
					i++
					if i > formatsLen {
						i = formatsLen
					}
					currFormat = formats[i]
				}
			} else {
				i++
				if i > formatsLen {
					i = formatsLen
				}
				currFormat = formats[i]
			}
		}
	}
	return
}

type StringHint struct {
	Repeat     bool
	Length     int
	Format     string
	Separator  byte
	Terminator byte
	Numeric    bool
}

const (
	StateRepeat = iota
	StateLength
	StateFormat
	StateSeparator
	StateTerminator
	StateEnd
)

func parseHint(format string) (formats []StringHint) {
	formatLen := len(format)
	formats = make([]StringHint, 0)
	state := StateRepeat
	var (
		currFormat StringHint
		c          byte
		i          int
	)
	for i < formatLen {
		c = format[i]
		switch state {
		case StateRepeat:
			if c == '*' {
				currFormat.Repeat = true
				i++
				break
			}
			state = StateLength
		case StateLength:
			if c >= '0' && c <= '9' {
				currFormat.Length = 10*currFormat.Length + int(c-'0')
				i++
				break
			}
			state = StateFormat
		case StateFormat:
			// TODO: Deal with UTF-8 (t) correctly
			if c == 'a' || c == 't' {
				c = 'c'
			} else {
				currFormat.Numeric = true
			}
			if c == 'x' {
				currFormat.Format = "%02x"
			} else {
				currFormat.Format = string([]byte{'%', c})
			}
			i++
			state = StateSeparator
		case StateSeparator:
			if c == '*' || (c >= '0' && c <= '9') {
				state = StateEnd
				break
			}
			currFormat.Separator = c
			i++
			if currFormat.Repeat {
				state = StateTerminator
			} else {
				state = StateEnd
			}
		case StateTerminator:
			if c == '*' || (c >= '0' && c <= '9') {
				state = StateEnd
				break
			}
			currFormat.Terminator = c
			i++
			state = StateEnd
		case StateEnd:
			formats = append(formats, currFormat)
			currFormat = StringHint{}
			state = StateRepeat
		}
	}
	formats = append(formats, currFormat)
	return
}

func pluralizeCount(i int64, base string) string {
	s := strconv.FormatInt(i, 10) + " " + base
	if i > 1 {
		s += "s"
	}
	return s
}

func DurationFormat(d time.Duration) string {
	seconds := int64(d / time.Second)
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	format := ""

	if days > 0 {
		format = strconv.FormatInt(days, 10) + "d "
	}
	if len(format) > 0 || hours%24 > 0 {
		format += strconv.FormatInt(hours%24, 10) + "h "
	}
	if len(format) > 0 || minutes%60 > 0 {
		format += strconv.FormatInt(minutes%60, 10) + "m "
	}
	if hours == 0 && seconds%60 > 0 {
		format += strconv.FormatInt(seconds%60, 10) + "s"
	}
	if format == "" {
		return "0s"
	}
	return strings.TrimSpace(format)
}

func DurationFormatLong(d time.Duration) string {
	seconds := int64(d / time.Second)
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	format := ""

	if days > 0 {
		format = pluralizeCount(days, "day") + " "
	}
	if hours%24 > 0 {
		format += pluralizeCount(hours%24, "hour") + " "
	}
	if minutes%60 > 0 {
		format += pluralizeCount(minutes%60, "min") + " "
	}
	if seconds%60 > 0 {
		format += pluralizeCount(seconds%60, "sec")
	}
	if format == "" {
		return "0 secs"
	}
	return strings.TrimSpace(format)
}
