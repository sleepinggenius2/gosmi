package models

import (
	"strconv"
	"strings"
	"time"
)

func GetDurationFormatted(value interface{}, flags Format) (v Value) {
	intVal, err := ToInt64(value)
	if err != nil {
		return
	}
	duration := time.Duration(intVal * 1e7)
	v.Format = flags
	v.Raw = duration
	if flags == FormatNone {
		return
	}
	if flags&FormatDurationShort > 0 {
		v.Formatted = DurationFormat(duration)
	} else {
		v.Formatted = DurationFormatLong(duration)
	}
	return
}

func GetDurationFormatter(flags Format) (f ValueFormatter) {
	return func(value interface{}) Value {
		return GetDurationFormatted(value, flags)
	}
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
