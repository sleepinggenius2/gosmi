package parser

import (
	"strconv"

	"github.com/alecthomas/participle/lexer"
	"github.com/pkg/errors"

	"github.com/sleepinggenius2/gosmi/types"
)

type SubIdentifier struct {
	Pos    lexer.Position
	Name   *types.SmiIdentifier
	Number *types.SmiSubId
}

func (x *SubIdentifier) Parse(lex *lexer.PeekingLexer) error {
	token, err := lex.Next()
	if err != nil {
		return err
	}
	x.Pos = token.Pos
	symbols := smiLexer.Symbols()
	if token.Type == symbols["Int"] {
		n, err := strconv.ParseUint(token.Value, 10, 32)
		if err != nil {
			return errors.Wrap(err, "Parse number")
		}
		x.Number = new(types.SmiSubId)
		*x.Number = types.SmiSubId(n)
		return nil
	} else if token.Type != symbols["Ident"] {
		return errors.Errorf("Unexpected %q, expected Ident", token)
	}
	x.Name = new(types.SmiIdentifier)
	*x.Name = types.SmiIdentifier(token.Value)
	token, err = lex.Peek(0)
	if err != nil {
		return err
	}
	if token.Value != "(" {
		return nil
	}
	_, err = lex.Next()
	if err != nil {
		return err
	}
	token, err = lex.Next()
	if err != nil {
		return err
	}
	if token.Type != symbols["Int"] {
		return errors.Errorf("Unexpected %q, expected Int", token)
	}
	n, err := strconv.ParseUint(token.Value, 10, 32)
	if err != nil {
		return errors.Wrap(err, "Parse number")
	}
	x.Number = new(types.SmiSubId)
	*x.Number = types.SmiSubId(n)
	token, err = lex.Next()
	if err != nil {
		return err
	}
	if token.Value != ")" {
		return errors.Errorf("Unexpected %q, expected \")\"", token)
	}
	return nil
}

type Oid struct {
	Pos lexer.Position

	SubIdentifiers []SubIdentifier `parser:"@@+"`
}

type Range struct {
	Pos lexer.Position

	Start string `parser:"@( \"-\"? Int | BinString | HexString )"`
	End   string `parser:"( \"..\" @( \"-\"? Int | BinString | HexString ) )?"`
}

type Status string

const (
	StatusMandatory  Status = "mandatory"
	StatusOptional   Status = "optional"
	StatusCurrent    Status = "current"
	StatusDeprecated Status = "deprecated"
	StatusObsolete   Status = "obsolete"
)

func (s Status) ToSmi() types.Status {
	switch s {
	case StatusMandatory:
		return types.StatusMandatory
	case StatusOptional:
		return types.StatusOptional
	case StatusCurrent:
		return types.StatusCurrent
	case StatusDeprecated:
		return types.StatusDeprecated
	case StatusObsolete:
		return types.StatusObsolete
	}
	return types.StatusUnknown
}

type SubType struct {
	Pos lexer.Position

	OctetString []Range `parser:"( ( \"SIZE\" \"(\" ( @@ ( \"|\" @@ )* ) \")\" )"`
	Integer     []Range `parser:"| @@ ( \"|\" @@ )* )"`
}

type NamedNumber struct {
	Pos lexer.Position

	Name  types.SmiIdentifier `parser:"@Ident"`
	Value string              `parser:"\"(\" @( \"-\"? Int ) \")\""`
}

type SyntaxType struct {
	Pos lexer.Position

	Name    types.SmiIdentifier `parser:"@( OctetString | ObjectIdentifier | Ident )"`
	SubType *SubType            `parser:"( ( \"(\" @@ \")\" )"`
	Enum    []NamedNumber       `parser:"| ( \"{\" @@ ( \",\" @@ )* \",\"? \"}\" ) )?"`
}

type Syntax struct {
	Pos lexer.Position

	Sequence *types.SmiIdentifier `parser:"( \"SEQUENCE\" \"OF\" @Ident )"`
	Type     *SyntaxType          `parser:"| @@"`
}
