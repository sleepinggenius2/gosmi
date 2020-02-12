package parser

import (
	"github.com/alecthomas/participle/lexer"
	"github.com/pkg/errors"

	"github.com/sleepinggenius2/gosmi/types"
)

type MacroBody struct {
	Pos lexer.Position

	TypeNotation  string
	ValueNotation string
	Tokens        map[string]string
}

func (m *MacroBody) Parse(lex *lexer.PeekingLexer) error {
	token, err := lex.Next()
	if err != nil {
		return err
	}
	if token.Value != "BEGIN" {
		return errors.Errorf("Expected 'BEGIN', Got '%s'", token.Value)
	}
	m.Pos = token.Pos

	var tokenName, tokenValue string
	m.Tokens = make(map[string]string)
	symbols := smiLexer.Symbols()
	for {
		token, err = lex.Next()
		if err != nil {
			return err
		}
		if token.Value == "END" {
			break
		}
		peek, _ := lex.Peek(0)
		if ((token.Value == "TYPE" || token.Value == "VALUE") && peek.Value == "NOTATION") || peek.Type == symbols["Assign"] {
			if token.Value == "NOTATION" {
				tokenName += " NOTATION"
				_, err = lex.Next()
				if err != nil {
					return err
				}
				continue
			}
			if tokenName != "" {
				switch tokenName {
				case "TYPE NOTATION":
					m.TypeNotation = tokenValue
				case "VALUE NOTATION":
					m.ValueNotation = tokenValue
				default:
					m.Tokens[tokenName] = tokenValue
				}
			}
			tokenName = token.Value
			tokenValue = ""
			if peek.Type == symbols["Assign"] {
				_, err = lex.Next()
				if err != nil {
					return err
				}
			}
			continue
		}
		if len(tokenValue) > 0 {
			tokenValue += " "
		}
		if token.Type == symbols["Text"] {
			tokenValue += `"` + token.Value + `"`
		} else {
			tokenValue += token.Value
		}
	}
	switch tokenName {
	case "":
		break
	case "TYPE NOTATION":
		m.TypeNotation = tokenValue
	case "VALUE NOTATION":
		m.ValueNotation = tokenValue
	default:
		m.Tokens[tokenName] = tokenValue
	}
	return nil
}

type Macro struct {
	Pos lexer.Position

	Name types.SmiIdentifier `parser:"@Ident \"MACRO\" Assign"`
	Body MacroBody           `parser:"@@"`
}
