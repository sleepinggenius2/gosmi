package parser

import (
	"github.com/alecthomas/participle/lexer"
)

type NotificationGroup struct {
	Pos lexer.Position

	Notifications []Identifier `parser:"\"NOTIFICATIONS\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\""` // Required
	Status        Status       `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"`   // Required
	Description   string       `parser:"\"DESCRIPTION\" @Text"`                                         // Required
	Reference     *string      `parser:"( \"REFERENCE\" @Text )?"`
	Oid           Oid          `parser:"Assign \"{\" @@ \"}\""`
}

type NotificationType struct {
	Pos lexer.Position

	Objects     []Identifier `parser:"( \"OBJECTS\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\" )?"`
	Status      Status       `parser:"\"STATUS\" @( \"current\" | \"deprecated\" | \"obsolete\" )"` // Required
	Description string       `parser:"\"DESCRIPTION\" @Text"`                                       // Required
	Reference   *string      `parser:"( \"REFERENCE\" @Text )?"`
	Oid         Oid          `parser:"Assign \"{\" @@ \"}\""`
}

type TrapType struct {
	Pos lexer.Position

	Enterprise    Identifier   `parser:"\"ENTERPRISE\" @Ident"`
	Objects       []Identifier `parser:"( \"VARIABLES\" \"{\" @Ident ( \",\" @Ident )* \",\"? \"}\" )?"`
	Description   string       `parser:"( \"DESCRIPTION\" @Text )?"`
	Reference     *string      `parser:"( \"REFERENCE\" @Text )?"`
	SubIdentifier uint32       `parser:"Assign @Int"`
}
