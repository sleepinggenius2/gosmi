package parser

import (
	"io"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

var (
	// TODO: Comments can also end with a "--"
	smiLexer = lexer.Must(ebnf.New(`
		Keyword = "FROM" .
		ObjectIdentifier = "OBJECT IDENTIFIER" .
		OctetString = "OCTET STRING" .
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
	smiParser = participle.MustBuild(new(Module),
		participle.Lexer(smiLexer),
		participle.Map(func(token lexer.Token) (lexer.Token, error) {
			if token.EOF() {
				return token, nil
			}
			token.Value = strings.Trim(token.Value, `"`)
			return token, nil
		}, "ExtUTCTime", "Text"),
		//participle.UseLookahead(2),
		participle.Upper("ExtUTCTime"),
		participle.Elide("Whitespace", "Comment"),
	)
)

func Parse(r io.Reader) (*Module, error) {
	m := new(Module)
	return m, smiParser.Parse(r, m)
}
