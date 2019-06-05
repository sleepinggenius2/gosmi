package parser

import (
	"github.com/alecthomas/participle/lexer"

	"github.com/sleepinggenius2/gosmi/types"
)

type TextualConvention struct {
	Pos lexer.Position

	DisplayHint string     `parser:"( \"DISPLAY-HINT\" @Text )?"`
	Status      Status     `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string     `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   string     `parser:"( \"REFERENCE\" @Text )?"`
	Syntax      SyntaxType `parser:"\"SYNTAX\" @@"` // Required
}

type SequenceEntry struct {
	Pos lexer.Position

	Descriptor types.SmiIdentifier `parser:"@Ident"`
	Syntax     SyntaxType          `parser:"@@"`
}

type SequenceType string

const (
	SequenceTypeChoice   SequenceType = "CHOICE"
	SequenceTypeSequence SequenceType = "SEQUENCE"
)

type Sequence struct {
	Pos lexer.Position

	Type    SequenceType    `parser:"@( \"CHOICE\" | \"SEQUENCE\" )"`
	Entries []SequenceEntry `parser:"\"{\" @@ ( \",\" @@ )* \",\"? \"}\""`
}

type Implicit struct {
	Pos lexer.Position

	Application bool       `parser:"\"[\" @\"APPLICATION\"?"`
	Number      int        `parser:"@Int \"]\""`
	Syntax      SyntaxType `parser:"\"IMPLICIT\" @@"`
}

type Type struct {
	Pos lexer.Position

	Name              types.SmiIdentifier `parser:"@Ident Assign"`
	TextualConvention *TextualConvention  `parser:"( ( \"TEXTUAL-CONVENTION\" @@ )"`
	Sequence          *Sequence           `parser:"| @@"`
	Implicit          *Implicit           `parser:"| @@"`
	Syntax            *SyntaxType         `parser:"| @@ )"`
}
