package parser

import (
	"time"

	"github.com/alecthomas/participle/lexer"

	"github.com/sleepinggenius2/gosmi/types"
)

type Date string

func (d Date) ToTime() (t time.Time) {
	if len(d) == 11 {
		t, _ = time.Parse("0601021504Z", string(d))
	} else {
		t, _ = time.Parse("200601021504Z", string(d))
	}
	return
}

type Import struct {
	Pos lexer.Position

	Names  []types.SmiIdentifier `parser:"@Ident ( \",\" @Ident )* \",\"?"`
	Module types.SmiIdentifier   `parser:"\"FROM\" @Ident"`
}

type Node struct {
	Pos lexer.Position

	Name              types.SmiIdentifier `parser:"@Ident"`
	ObjectIdentifier  bool                `parser:"( ( ( @ObjectIdentifier"`
	ObjectIdentity    *ObjectIdentity     `parser:"| ( \"OBJECT-IDENTITY\" @@ )"`
	ObjectGroup       *ObjectGroup        `parser:"| ( \"OBJECT-GROUP\" @@ )"`
	ObjectType        *ObjectType         `parser:"| ( \"OBJECT-TYPE\" @@ )"`
	NotificationGroup *NotificationGroup  `parser:"| ( \"NOTIFICATION-GROUP\" @@ )"`
	NotificationType  *NotificationType   `parser:"| ( \"NOTIFICATION-TYPE\" @@ )"`
	ModuleCompliance  *ModuleCompliance   `parser:"| ( \"MODULE-COMPLIANCE\" @@ )"`
	AgentCapabilities *AgentCapabilities  `parser:"| ( \"AGENT-CAPABILITIES\" @@ ) )"`
	Oid               *Oid                `parser:"Assign \"{\" @@ \"}\" )"`
	TrapType          *TrapType           `parser:"| ( ( \"TRAP-TYPE\" @@ )"`
	SubIdentifier     *types.SmiSubId     `parser:"Assign @Int ) )"`
}

type Revision struct {
	Pos lexer.Position

	Date        Date   `parser:"@ExtUTCTime"`
	Description string `parser:"\"DESCRIPTION\" @Text"`
}

type ModuleIdentity struct {
	Pos lexer.Position

	Name         types.SmiIdentifier `parser:"@Ident \"MODULE-IDENTITY\""`
	LastUpdated  Date                `parser:"\"LAST-UPDATED\" @ExtUTCTime"` // Required
	Organization string              `parser:"\"ORGANIZATION\" @Text"`       // Required
	ContactInfo  string              `parser:"\"CONTACT-INFO\" @Text"`       // Required
	Description  string              `parser:"\"DESCRIPTION\" @Text"`        // Required
	Revisions    []Revision          `parser:"( \"REVISION\" @@ )*"`
	Oid          Oid                 `parser:"Assign \"{\" @@ \"}\""`
}

type ModuleBody struct {
	Pos lexer.Position

	Imports  []Import              `parser:"( \"IMPORTS\" @@+ \";\" )?"`
	Exports  []types.SmiIdentifier `parser:"( \"EXPORTS\"  @Ident ( \",\" @Ident )* \";\" )?"`
	Identity *ModuleIdentity       `parser:"( @@"`
	Types    []Type                `parser:"| @@"`
	Nodes    []Node                `parser:"| @@"`
	Macros   []Macro               `parser:"| @@ )*"`
}

type Module struct {
	Pos lexer.Position

	Name types.SmiIdentifier `parser:"@Ident"`
	Body ModuleBody          `parser:"\"DEFINITIONS\" Assign \"BEGIN\" @@ \"END\""`
}
