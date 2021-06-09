package parser

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

var (
	// TODO: Comments can also end with a "--"
	// Per the ASN.1 (ITU-T X.680) specification of a number token (Int below):
	// The first digit shall not be zero unless the "number" is a single digit.
	smiLexer = lexer.Must(ebnf.New(`
		Keyword = "FROM" .
		ObjectIdentifier = "OBJECT" Whitespace { Whitespace } "IDENTIFIER" .
		OctetString = "OCTET" Whitespace { Whitespace } "STRING" .
		BinString = "'" { "0" | "1" } "'" ( "b" | "B" ) .
		HexString = "'" { digit | "a"…"f" | "A"…"F" } "'" ( "h" | "H" ) .
		Assign = "::=" .
		Comment = "--" { "\u0000"…"\U0010ffff"-"\n" } .
		ExtUTCTime = "\"" digit digit digit digit digit digit digit digit digit digit [ digit digit ] ( "z" | "Z" ) "\"" .
		Text = "\"" { "\u0000"…"\U0010ffff"-"\"" } "\"" .
		Ident = alpha { alpha | digit | "-" | "_" } .
		Int = "0" | ( digit { digit } ) .
		Punct = ".." | "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
		Whitespace = " " | "\t" | "\n" | "\r" .

		lower = "a"…"z" .
		upper = "A"…"Z" .
		alpha = lower | upper .
		digit = "0"…"9" .
	`))
	compressSpace = regexp.MustCompile(`(?:\r?\n *)+`)
	smiParser     = participle.MustBuild(new(Module),
		participle.Lexer(smiLexer),
		participle.Map(func(token lexer.Token) (lexer.Token, error) {
			if token.EOF() {
				return token, nil
			}
			token.Value = compressSpace.ReplaceAllString(strings.TrimSpace(strings.Trim(token.Value, `"`)), "\n")
			return token, nil
		}, "ExtUTCTime", "Text"),
		participle.Map(func(token lexer.Token) (lexer.Token, error) {
			token.Value = "OBJECT IDENTIFIER"
			return token, nil
		}, "ObjectIdentifier"),
		participle.Map(func(token lexer.Token) (lexer.Token, error) {
			token.Value = "OCTET STRING"
			return token, nil
		}, "OctetString"),
		//participle.UseLookahead(2),
		participle.Upper("ExtUTCTime", "BinString", "HexString"),
		participle.Elide("Whitespace", "Comment"),
	)
)

func Parse(r io.Reader) (*Module, error) {
	m := new(Module)
	return m, smiParser.Parse(r, m)
}

func ParseFile(path string) (*Module, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Open file: %w", err)
	}
	defer r.Close()
	module, err := Parse(r)
	if err != nil {
		return module, fmt.Errorf("Parse file: %w", err)
	}
	return module, nil
}
