package parser

import (
	"github.com/alecthomas/participle/lexer"
)

type ObjectGroup struct {
	Pos lexer.Position

	Objects     []Identifier `parser:"\"OBJECTS\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\""`     // Required
	Status      Status       `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string       `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   *string      `parser:"( \"REFERENCE\" @Text )?"`
	Oid         Oid          `parser:"Assign \"{\" @@ \"}\""`
}

type ObjectIdentifier struct {
	Pos lexer.Position

	Oid Oid `parser:"Assign \"{\" @@ \"}\""`
}

type ObjectIdentity struct {
	Pos lexer.Position

	Status      Status  `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string  `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   *string `parser:"( \"REFERENCE\" @Text )?"`
	Oid         Oid     `parser:"Assign \"{\" @@ \"}\""`
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

type Index struct {
	Pos lexer.Position

	Implied bool       `parser:"@\"IMPLIED\"?"`
	Name    Identifier `parser:"@Ident"`
}

type ObjectType struct {
	Pos lexer.Position

	Syntax      Syntax      `parser:"\"SYNTAX\" @@"` // Required
	Units       *string     `parser:"( \"UNITS\" @Text )?"`
	MaxAccess   Access      `parser:"( \"ACCESS\" | \"MAX-ACCESS\" ) @( \"write-only\" | \"not-accessible\" | \"accessible-for-notify\" | \"read-only\" | \"read-write\" | \"read-create\" )"` // Required
	Status      Status      `parser:"\"STATUS\" @( \"mandatory\" | \"optional\" | \"current\" | \"deprecated\" | \"obsolete\" )"`                                                              // Required
	Description string      `parser:"( \"DESCRIPTION\" @Text )?"`                                                                                                                              // Required RFC 1212+
	Reference   *string     `parser:"( \"REFERENCE\" @Text )?"`
	Index       []Index     `parser:"( ( \"INDEX\" \"{\" @@ ( \",\" @@ )* \",\"? \"}\" )"` // Required for "row" without AUGMENTS
	Augments    *Identifier `parser:"| ( \"AUGMENTS\" \"{\" @Ident \"}\" ) )?"`            // Required for "row" without INDEX
	Defval      *string     `parser:"( \"DEFVAL\" \"{\" @( \"-\"? Int | BinString | HexString | Text | Ident | ( \"{\" ( Int+ | ( Ident ( \",\" Ident )* \",\"? )? ) \"}\" ) ) \"}\" )?"`
	Oid         Oid         `parser:"Assign \"{\" @@ \"}\""`
}
