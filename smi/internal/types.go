package internal

import (
	"github.com/sleepinggenius2/gosmi/types"
)

type List struct {
	types.SmiElement
	Ptr  interface{}
	Next *List
}

type Kind int

const (
	KindUnknown  Kind = 0
	KindModule   Kind = 1
	KindMacro    Kind = 2
	KindType     Kind = 3
	KindObject   Kind = 4
	KindImport   Kind = 5
	KindImported Kind = 6
	KindNotFound Kind = 7
)

type Flags uint16

const (
	FlagRoot         Flags = 0x0001 // Mark node tree's root
	FlagSeqType      Flags = 0x0002 // Type is set from SMIv1/2 SEQUENCE
	FlagRegistered   Flags = 0x0004 // On an Object: this is registered
	FlagIncomplete   Flags = 0x0008 // Just defined by a forward referenced type or object
	FlagCreatable    Flags = 0x0040 // On a Row: new rows can be created
	FlagInGroup      Flags = 0x0080 // Node is contained in a group
	FlagInCompliance Flags = 0x0100 // Group is mentioned in a compliance statement. In case of ImportFlags: the import is done through a compliance MODULE phrase
	FlagInSyntax     Flags = 0x0200 // Type is mentioned in a syntax statement
)

func (x Flags) Has(flag Flags) bool {
	return x&flag == flag
}

const (
	UnknownLabel = "<unknown>"
)
