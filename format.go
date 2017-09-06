package gosmi

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
	FormatAll      Format = 0xff & ^FormatUnits
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

type ValueFormatter func(interface{}) Value

func (n Node) GetValueFormatter(flags ...Format) (f ValueFormatter) {
	formatFlags := FormatNone
	for _, flag := range flags {
		formatFlags |= flag
	}
	if n.Type.BaseType == types.BaseTypeOctetString {
		f = func(value interface{}) (v Value) {
			v.Format = formatFlags
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
				if formatFlags & FormatString != 0 {
					v.Formatted = fmt.Sprintf("%v", val)
				}
				return
			}
			v.Raw = bytes
			if formatFlags & FormatString == 0 {
				return
			}
			format := n.Type.Format
			if n.Type.Name == "IpAddress" {
				format = "1d."
			} else if n.Type.Name == "InetAddress" {
				if bytes[0] == 0x04 {
					format = "*1d."
				} else {
					format = "*1x:"
				}
			} else if n.Type.Name == "IpV4orV6Addr" {
				if len(bytes) == 4 {
					format = "1d."
				} else {
					format = "1x:"
				}
			}
			formatted := StringDisplayHint(format, bytes)
			if formatFlags & FormatUnits != 0 && n.Type.Units != "" {
				formatted += " " + n.Type.Units
			}
			v.Formatted = formatted
			return
		}
		return
	} else if n.Type.Name == "TimeTicks" || n.Type.Name == "TimeInterval" || n.Type.Name == "TimeStamp" {
		f = func(value interface{}) (v Value) {
			v.Format = formatFlags
			duration := time.Duration(int64(value.(int)) * 1e7)
			v.Raw = duration
			v.Formatted = DurationFormat(duration)
			return
		}
		return
	}
	if n.Type.Enum == nil {
		f = func(value interface{}) (v Value) {
			v.Format = formatFlags
			v.Raw = value
			if formatFlags == FormatNone {
				return
			}
			if n.Type.BaseType == types.BaseTypeBits {
				if formatFlags & FormatBits != 0 {
					v.Formatted = fmt.Sprintf("% X", value.([]byte))
				}
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
			formatted := IntegerDisplayHint(n.Type.Format, intVal)
			if formatFlags & FormatUnits != 0 && n.Type.Units != "" {
				formatted += " " + n.Type.Units
			}
			v.Formatted = formatted
			return
		}
		return
	}
	var enums map[int64]string
	enums = make(map[int64]string, len(n.Type.Enum.Values))
	for _, enum := range n.Type.Enum.Values {
		intVal := gosnmp.ToBigInt(enum.Value).Int64()
		enums[intVal] = enum.Name
	}
	f = func(value interface{}) (v Value) {
		v.Format = formatFlags
		v.Raw = value
		if formatFlags == FormatNone {
			return
		}
		if n.Type.BaseType == types.BaseTypeBits {
			octets := value.([]byte)
			if formatFlags & FormatBits != 0 {
				v.Formatted = fmt.Sprintf("% X", octets)
			}
			if (formatFlags & FormatEnumName) + (formatFlags & FormatEnumValue) == 0 {
				return
			}
			bitsFormatted := make([]string, 0, 8*len(octets))
			for i, octet := range octets {
				for j := 7; j >= 0; j-- {
					if octet&(1<<uint(j)) != 0 {
						bit := 8*i + (7 - j)
						var bitFormatted string
						if formatFlags & FormatEnumName != 0 {
							enumName, ok := enums[int64(bit)]
							if !ok {
								enumName = "unknown"
							}
							bitFormatted = enumName
							if formatFlags & FormatEnumValue != 0 || !ok {
								bitFormatted += "(" + fmt.Sprintf("%d", bit) + ")"
							}
						} else if formatFlags & FormatEnumValue != 0 {
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
		} else {
			var intVal int64
			switch tempVal := value.(type) {
			case int64:
				intVal = tempVal
			default:
				intVal = gosnmp.ToBigInt(tempVal).Int64()
			}
			v.Raw = intVal
			if formatFlags & FormatEnumName != 0 {
				enumName, ok := enums[intVal]
				if !ok {
					enumName = "unknown"
				}
				v.Formatted = enumName
				if formatFlags & FormatEnumValue != 0 {
					v.Formatted += fmt.Sprintf("(%d)", intVal)
				}
			} else if formatFlags & FormatEnumValue != 0 {
				v.Formatted = fmt.Sprintf("%d", intVal)
			}
		}
		return
	}
	return
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
		if formattedLen - offset <= decimals {
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
	var i, repeatCount, lengthCount int
	var currValue uint64
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
			currValue = currValue<<8+uint64(c)
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
		c byte
		i int
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
	seconds := int64(d/time.Second)
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	format := ""

	if days > 0 {
		format = strconv.FormatInt(days, 10) + "d "
	}
	if len(format) > 0 || hours % 24 > 0 {
		format += strconv.FormatInt(hours % 24, 10) + "h "
	}
	if len(format) > 0 || minutes % 60 > 0 {
		format += strconv.FormatInt(minutes % 60, 10) + "m "
	}
	if hours == 0 && seconds % 60 > 0 {
		format += strconv.FormatInt(seconds % 60, 10) + "s"
	}
	if format == "" {
		return "0s"
	}
	return strings.TrimSpace(format)
}

func DurationFormatLong(d time.Duration) string {
	seconds := int64(d/time.Second)
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	format := ""

	if days > 0 {
		format = pluralizeCount(days, "day") + " "
	}
	if hours % 24 > 0 {
		format += pluralizeCount(hours % 24, "hour") + " "
	}
	if minutes % 60 > 0 {
		format += pluralizeCount(minutes % 60, "min") + " "
	}
	if seconds % 60 > 0 {
		format += pluralizeCount(seconds % 60, "sec")
	}
	if format == "" {
		return "0 secs"
	}
	return strings.TrimSpace(format)
}
