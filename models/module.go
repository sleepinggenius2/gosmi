package models

import (
	"github.com/sleepinggenius2/gosmi/types"
)

type Module struct {
	ContactInfo  string
	Description  string
	Language     types.Language
	Name         string
	Organization string
	Path         string
	Reference    string
}
