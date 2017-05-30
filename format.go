package gosmi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sleepinggenius2/gosmi/types"
	"github.com/sleepinggenius2/gosnmp"
)

type ValueFormatter func(interface{}) string

func (n Node) GetValueFormatter() (f ValueFormatter) {
	if n.Type.BaseType == types.BaseTypeOctetString {
		f = func(value interface{}) string {
			var bytes []byte
			format := n.Type.Format
			switch val := value.(type) {
			case []int:
				bytes = make([]byte, len(val))
				for i, b := range val {
					bytes[i] = byte(b)
				}
			case []byte:
				bytes = val
			default:
				return fmt.Sprintf("%v", val)
			}
			if n.Type.Name == "IpAddress" {
				format = "1d."
			} else if n.Type.Name == "InetAddress" {
				if bytes[0] == 4 {
					format = "*1d."
				} else {
					format = "*1x:"
				}
			}
			formatted := StringDisplayHint(format, bytes)
			if n.Type.Units != "" {
				formatted += " " + n.Type.Units
			}
			return formatted
		}
		return
	}
	if n.Type.Enum == nil {
		f = func(value interface{}) string {
			if n.Type.BaseType != types.BaseTypeBits {
				formatted := IntegerDisplayHint(n.Type.Format, value.(int64))
				if n.Type.Units != "" {
					formatted += " " + n.Type.Units
				}
				return formatted
			}
			return fmt.Sprintf("% X", value.([]byte))
		}
		return
	}
	var enums map[int64]string
	enums = make(map[int64]string, len(n.Type.Enum.Values))
	for _, enum := range n.Type.Enum.Values {
		intVal := gosnmp.ToBigInt(enum.Value).Int64()
		enums[intVal] = enum.Name
	}
	f = func(value interface{}) (formattedValue string) {
		if n.Type.BaseType == types.BaseTypeBits {
			octets := value.([]byte)
			octetsFormatted := make([]string, len(octets))
			bitsFormatted := make([]string, 0, 8*len(octets))
			for i, octet := range octets {
				octetsFormatted[i] = fmt.Sprintf("%02X", octet)
				for j := 7; j >= 0; j-- {
					if octet&(1<<uint(j)) != 0 {
						bit := 8*i + (7 - j)
						bitFormatted := fmt.Sprintf("%d", bit)
						if enums != nil {
							val, ok := enums[int64(bit)]
							if ok {
								bitFormatted = val + "(" + bitFormatted + ")"
							}
						}
						bitsFormatted = append(bitsFormatted, bitFormatted)
					}
				}
			}
			formattedValue = fmt.Sprintf("%s[%s]", strings.Join(octetsFormatted, " "), strings.Join(bitsFormatted, " "))
		} else {
			intVal := value.(int64)
			name, ok := enums[intVal]
			if !ok {
				name = "unknown"
			}
			formattedValue = fmt.Sprintf("%s(%d)", name, intVal)
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
		if formattedLen <= decimals {
			formatStr := "0.%0" + format[2:] + "s"
			formatted = fmt.Sprintf(formatStr, formatted)
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
		return string(value)
	}
	formats := parseHint(format)
	formatsLen := len(formats) - 1
	//fmt.Printf("%s = %+v\n", format, formats)
	var i, repeatCount, lengthCount int
	for j, c := range value {
		if i > formatsLen {
			i = formatsLen
		}
		if formats[i].Repeat && repeatCount == 0 {
			repeatCount = int(c)
			continue
		}
		if lengthCount == 0 {
			lengthCount = formats[i].Length
		}
		lengthCount--
		formatted += fmt.Sprintf(formats[i].Format, c)
		if lengthCount == 0 {
			if formats[i].Separator > 0 && j != len(value)-1 {
				formatted += string([]byte{formats[i].Separator})
			}
			if formats[i].Repeat {
				repeatCount--
				if repeatCount == 0 {
					if formats[i].Terminator > 0 && j != len(value)-1 {
						formatted += string([]byte{formats[i].Terminator})
					}
					i++
				}
			} else {
				i++
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
	var currFormat StringHint
	i := 0
	var c byte
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
		}
	}
	formats = append(formats, currFormat)
	return
}
