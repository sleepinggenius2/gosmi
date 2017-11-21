package models

import (
	"time"

	"github.com/sleepinggenius2/gosmi/types"
)

type Import struct {
	Module string
	Name   string
}

type Module struct {
	ContactInfo  string
	Description  string
	Language     types.Language
	Name         string
	Organization string
	Path         string
	Reference    string
}

type Revision struct {
	Date        time.Time
	Description string
}
