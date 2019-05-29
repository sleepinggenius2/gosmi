package parser

import (
	"github.com/alecthomas/participle/lexer"
)

type Import struct {
	Pos lexer.Position

	Names  []Identifier `parser:"@Ident ( \",\" @Ident )* \",\"?"`
	Module Identifier   `parser:"\"FROM\" @Ident"`
}

type Node struct {
	Pos lexer.Position

	Name              Identifier         `parser:"@Ident"`
	ObjectIdentifier  *ObjectIdentifier  `parser:"( ( ObjectIdentifier @@ )"`
	ObjectIdentity    *ObjectIdentity    `parser:"| ( \"OBJECT-IDENTITY\" @@ )"`
	ObjectGroup       *ObjectGroup       `parser:"| ( \"OBJECT-GROUP\" @@ )"`
	ObjectType        *ObjectType        `parser:"| ( \"OBJECT-TYPE\" @@ )"`
	NotificationGroup *NotificationGroup `parser:"| ( \"NOTIFICATION-GROUP\" @@ )"`
	NotificationType  *NotificationType  `parser:"| ( \"NOTIFICATION-TYPE\" @@ )"`
	TrapType          *TrapType          `parser:"| ( \"TRAP-TYPE\" @@ )"`
	ModuleCompliance  *ModuleCompliance  `parser:"| ( \"MODULE-COMPLIANCE\" @@ )"`
	AgentCapabilities *AgentCapabilities `parser:"| ( \"AGENT-CAPABILITIES\" @@ ) )"`
}

type Revision struct {
	Pos lexer.Position

	Date        string `parser:"@ExtUTCTime"`
	Description string `parser:"\"DESCRIPTION\" @Text"`
}

type ModuleIdentity struct {
	Pos lexer.Position

	Name         Identifier `parser:"@Ident \"MODULE-IDENTITY\""`
	LastUpdated  string     `parser:"\"LAST-UPDATED\" @ExtUTCTime"` // Required
	Organization string     `parser:"\"ORGANIZATION\" @Text"`       // Required
	ContactInfo  string     `parser:"\"CONTACT-INFO\" @Text"`       // Required
	Description  string     `parser:"\"DESCRIPTION\" @Text"`        // Required
	Revisions    []Revision `parser:"( \"REVISION\" @@ )*"`
	Oid          Oid        `parser:"Assign \"{\" @@ \"}\""`
}

type ModuleBody struct {
	Pos lexer.Position

	Imports  []Import        `parser:"( \"IMPORTS\" @@+ \";\" )?"`
	Exports  []Identifier    `parser:"( \"EXPORTS\"  @Ident ( \",\" @Ident )* \";\" )?"`
	Identity *ModuleIdentity `parser:"( @@"`
	Types    []Type          `parser:"| @@"`
	Nodes    []Node          `parser:"| @@"`
	Macros   []Macro         `parser:"| @@ )*"`
}

type Module struct {
	Pos lexer.Position

	Name Identifier `parser:"@Ident"`
	Body ModuleBody `parser:"\"DEFINITIONS\" Assign \"BEGIN\" @@ \"END\""`
}
