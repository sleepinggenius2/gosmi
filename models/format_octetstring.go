package models

import (
	"fmt"
)

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

func GetInetAddressFormatted(value interface{}, flags Format) (v Value) {
	v.Format = flags
	if str, ok := value.(string); ok {
		v.Raw = value
		v.Formatted = str
		return
	}
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

func GetInetAddressFormatter(flags Format) (f ValueFormatter) {
	return func(value interface{}) Value {
		return GetInetAddressFormatted(value, flags)
	}
}

func GetOctetStringFormatted(value interface{}, flags Format, format string) (v Value) {
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

func GetOctetStringFormatter(flags Format, format string) (f ValueFormatter) {
	return func(value interface{}) Value {
		return GetOctetStringFormatted(value, flags, format)
	}
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
