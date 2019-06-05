package parser

import (
	"github.com/alecthomas/participle/lexer"

	"github.com/sleepinggenius2/gosmi/types"
)

type ObjectGroup struct {
	Pos lexer.Position

	Objects     []types.SmiIdentifier `parser:"\"OBJECTS\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\""`     // Required
	Status      Status                `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string                `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   string                `parser:"( \"REFERENCE\" @Text )?"`
}

type ObjectIdentity struct {
	Pos lexer.Position

	Status      Status `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   string `parser:"( \"REFERENCE\" @Text )?"`
}

type Access string

const (
	// In order from least to greatest
	AccessWriteOnly           Access = "write-only" // Do not use
	AccessNotImplemented      Access = "not-implemented"
	AccessNotAccessible       Access = "not-accessible"
	AccessAccessibleForNotify Access = "accesible-for-notify"
	AccessReadOnly            Access = "read-only"
	AccessReadWrite           Access = "read-write"
	AccessReadCreate          Access = "read-create"
)

func (a Access) ToSmi() types.Access {
	switch a {
	case AccessWriteOnly:
		return types.AccessUnknown // What should this be?
	case AccessNotImplemented:
		return types.AccessNotImplemented
	case AccessNotAccessible:
		return types.AccessNotAccessible
	case AccessAccessibleForNotify:
		return types.AccessNotify
	case AccessReadOnly:
		return types.AccessReadOnly
	case AccessReadWrite, AccessReadCreate:
		return types.AccessReadWrite
	}
	return types.AccessUnknown
}

type Index struct {
	Pos lexer.Position

	Implied bool                `parser:"@\"IMPLIED\"?"`
	Name    types.SmiIdentifier `parser:"@Ident"`
}

type ObjectType struct {
	Pos lexer.Position

	Syntax      Syntax               `parser:"\"SYNTAX\" @@"` // Required
	Units       string               `parser:"( \"UNITS\" @Text )?"`
	Access      Access               `parser:"( \"ACCESS\" | \"MAX-ACCESS\" ) @( \"write-only\" | \"not-accessible\" | \"accessible-for-notify\" | \"read-only\" | \"read-write\" | \"read-create\" )"` // Required
	Status      Status               `parser:"\"STATUS\" @( \"mandatory\" | \"optional\" | \"current\" | \"deprecated\" | \"obsolete\" )"`                                                              // Required
	Description string               `parser:"( \"DESCRIPTION\" @Text )?"`                                                                                                                              // Required RFC 1212+
	Reference   string               `parser:"( \"REFERENCE\" @Text )?"`
	Index       []Index              `parser:"( ( \"INDEX\" \"{\" @@ ( \",\" @@ )* \",\"? \"}\" )"` // Required for "row" without AUGMENTS
	Augments    *types.SmiIdentifier `parser:"| ( \"AUGMENTS\" \"{\" @Ident \"}\" ) )?"`            // Required for "row" without INDEX
	Defval      *string              `parser:"( \"DEFVAL\" \"{\" @( \"-\"? Int | BinString | HexString | Text | Ident | ( \"{\" ( Int+ | ( Ident ( \",\" Ident )* \",\"? )? ) \"}\" ) ) \"}\" )?"`
}
