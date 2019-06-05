package parser

import (
	"github.com/alecthomas/participle/lexer"

	"github.com/sleepinggenius2/gosmi/types"
)

type NotificationGroup struct {
	Pos lexer.Position

	Notifications []types.SmiIdentifier `parser:"\"NOTIFICATIONS\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\""` // Required
	Status        Status                `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"`   // Required
	Description   string                `parser:"\"DESCRIPTION\" @Text"`                                         // Required
	Reference     string                `parser:"( \"REFERENCE\" @Text )?"`
}

type NotificationType struct {
	Pos lexer.Position

	Objects     []types.SmiIdentifier `parser:"( \"OBJECTS\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\" )?"`
	Status      Status                `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string                `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   string                `parser:"( \"REFERENCE\" @Text )?"`
}

type TrapType struct {
	Pos lexer.Position

	Enterprise  types.SmiIdentifier   `parser:"\"ENTERPRISE\" @Ident"`
	Objects     []types.SmiIdentifier `parser:"( \"VARIABLES\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\" )?"`
	Description string                `parser:"( \"DESCRIPTION\" @Text )?"`
	Reference   string                `parser:"( \"REFERENCE\" @Text )?"`
}
