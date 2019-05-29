package parser

import (
	"github.com/alecthomas/participle/lexer"
)

type Identifier string

type SubIdentifier struct {
	Pos lexer.Position

	Name   *Identifier `parser:"@Ident?"`
	Number *uint32     `parser:"( \"(\" @Int \")\" | @Int )?"`
}

type Oid struct {
	Pos lexer.Position

	SubIdentifiers []SubIdentifier `parser:"@@+"`
}

type Range struct {
	Pos lexer.Position

	Start string  `parser:"@( \"-\"? Int | BinString | HexString )"`
	End   *string `parser:"( \"..\" @( \"-\"? Int | BinString | HexString ) )?"`
}

type Status string

const (
	StatusMandatory  Status = "mandatory"
	StatusOptional   Status = "optional"
	StatusCurrent    Status = "current"
	StatusDeprecated Status = "deprecated"
	StatusObsolete   Status = "obsolete"
)

type SubType struct {
	Pos lexer.Position

	OctetString []Range `parser:"( ( \"SIZE\" \"(\" ( @@ ( \"|\" @@ )* ) \")\" )"`
	Integer     []Range `parser:"| @@ ( \"|\" @@ )* )"`
}

type NamedNumber struct {
	Pos lexer.Position

	Name   Identifier `parser:"@Ident"`
	Number int64      `parser:"\"(\" @( \"-\"? Int ) \")\""`
}

type SyntaxType struct {
	Pos lexer.Position

	Name    Identifier    `parser:"@( OctetString | ObjectIdentifier | Ident )"`
	SubType *SubType      `parser:"( ( \"(\" @@ \")\" )"`
	Enum    []NamedNumber `parser:"| ( \"{\" @@ ( \",\" @@ )* \",\"? \"}\" ) )?"`
}

type Syntax struct {
	Pos lexer.Position

	Sequence *Identifier `parser:"( \"SEQUENCE\" \"OF\" @Ident )"`
	Type     *SyntaxType `parser:"| @@"`
}
