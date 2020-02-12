package parser

import (
	"github.com/alecthomas/participle/lexer"

	"github.com/sleepinggenius2/gosmi/types"
)

type AgentCapabilityVariation struct {
	Pos lexer.Position

	Name        types.SmiIdentifier   `parser:"\"VARIATION\" @Ident"` // Required
	Syntax      *Syntax               `parser:"( \"SYNTAX\" @@ )?"`
	WriteSyntax *Syntax               `parser:"( \"WRITE-SYNTAX\" @@ )?"`
	Access      *Access               `parser:"( \"ACCESS\" @( \"write-only\" | \"not-implemented\" | \"accessible-for-notify\" | \"read-only\" | \"read-write\" | \"read-create\" ) )?"`
	Creation    []types.SmiIdentifier `parser:"( \"CREATION-REQUIRES\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\" )?"`
	Defval      *string               `parser:"( \"DEFVAL\" \"{\" @( \"-\"? Int | BinString | HexString | Text | Ident | ( \"{\" ( Int+ | ( Ident ( \",\" Ident )* \",\"? )? ) \"}\" ) ) \"}\" )?"`
	Description string                `parser:"\"DESCRIPTION\" @Text"` // Required
}

type AgentCapabilityModule struct {
	Pos lexer.Position

	Module     types.SmiIdentifier        `parser:"\"SUPPORTS\" @Ident"`                                      // Required
	Includes   []types.SmiIdentifier      `parser:"\"INCLUDES\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\""` // Required
	Variations []AgentCapabilityVariation `parser:"@@*"`
}

type AgentCapabilities struct {
	Pos lexer.Position

	ProductRelease string                  `parser:"\"PRODUCT-RELEASE\" @Text"`                  // Required
	Status         Status                  `parser:"\"STATUS\" @( \"current\" | \"obsolete\" )"` // Required
	Description    string                  `parser:"\"DESCRIPTION\" @Text"`                      // Required
	Reference      string                  `parser:"( \"REFERENCE\" @Text )?"`
	Modules        []AgentCapabilityModule `parser:"@@*"`
}

type ComplianceGroup struct {
	Pos lexer.Position

	Name        types.SmiIdentifier `parser:"\"GROUP\" @Ident"`
	Description string              `parser:"\"DESCRIPTION\" @Text"`
}

type ComplianceObject struct {
	Pos lexer.Position

	Name        types.SmiIdentifier `parser:"\"OBJECT\" @Ident"`
	Syntax      *Syntax             `parser:"( \"SYNTAX\" @@ )?"`
	WriteSyntax *Syntax             `parser:"( \"WRITE-SYNTAX\" @@ )?"`
	MinAccess   *Access             `parser:"( \"MIN-ACCESS\" @( \"not-accessible\" | \"accessible-for-notify\" | \"read-only\" | \"read-write\" | \"read-create\" ) )?"`
	Description string              `parser:"\"DESCRIPTION\" @Text"`
}

type Compliance struct {
	Pos lexer.Position

	Group  *ComplianceGroup  `parser:"@@"`
	Object *ComplianceObject `parser:"| @@"`
}

type ComplianceModuleName string

func (n *ComplianceModuleName) Parse(lex *lexer.PeekingLexer) error {
	token, err := lex.Peek(0)
	if err != nil {
		return err
	}
	if token.Type == smiLexer.Symbols()["Assign"] || token.Value == "MANDATORY-GROUPS" || token.Value == "GROUP" || token.Value == "OBJECT" {
		*n = ""
		return nil
	}
	token, err = lex.Next()
	if err != nil {
		return err
	}
	*n = ComplianceModuleName(token.Value)
	return nil
}

type ModuleComplianceModule struct {
	Pos lexer.Position

	Name            ComplianceModuleName  `parser:"@@"`
	MandatoryGroups []types.SmiIdentifier `parser:"( \"MANDATORY-GROUPS\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\" )?"`
	Compliances     []Compliance          `parser:"@@*"`
}

type ModuleCompliance struct {
	Pos lexer.Position

	Status      Status                   `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string                   `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   string                   `parser:"( \"REFERENCE\" @Text )?"`
	Modules     []ModuleComplianceModule `parser:"( \"MODULE\" @@ )+"`
}
