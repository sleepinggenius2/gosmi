package types

import (
	"time"
)

type SmiModule struct {
	Name         SmiIdentifier
	Path         string
	Organization string
	ContactInfo  string
	Description  string
	Reference    string
	Language     Language
	Conformance  bool
}

type SmiRevision struct {
	Date        time.Time
	Description string
}

type SmiImport struct {
	Module SmiIdentifier
	Name   SmiIdentifier
}

type SmiValue struct {
	BaseType BaseType
	Len      uint
	Value    interface{}
}

type SmiType struct {
	Name        SmiIdentifier
	BaseType    BaseType
	Decl        Decl
	Format      string
	Value       SmiValue
	Units       string
	Status      Status
	Description string
	Reference   string
}

type SmiNamedNumber struct {
	Name  SmiIdentifier
	Value SmiValue
}

type SmiRange struct {
	MinValue SmiValue
	MaxValue SmiValue
}

type SmiNode struct {
	Name        SmiIdentifier
	OidLen      int
	Oid         Oid
	Decl        Decl
	Access      Access
	Status      Status
	Format      string
	Value       SmiValue
	Units       string
	Description string
	Reference   string
	IndexKind   IndexKind
	Implied     bool
	Create      bool
	NodeKind    NodeKind
}

type SmiElement struct{}

type SmiOption struct {
	Description string
}

type SmiRefinement struct {
	Access      Access
	Description string
}

type SmiMacro struct {
	Name        SmiIdentifier
	Decl        Decl
	Status      Status
	Description string
	Reference   string
}

// void (SmiErrorHandler) (char *path, int line, int severity, char *msg, char *tag)
type SmiErrorHandler func(path string, line int, severity int, msg string, tag string)
