package internal

import (
	"github.com/sleepinggenius2/gosmi/types"
)

type Macro struct {
	types.SmiMacro
	Module *Module
	Flags  Flags
	Next   *Macro
	Prev   *Macro
	Line   int
}

type MacroMap struct {
	First *Macro

	last *Macro
	m    map[types.SmiIdentifier]*Macro
}

func (x *MacroMap) Add(m *Macro) {
	m.Prev = x.last
	if x.First == nil {
		x.First = m
	} else {
		x.last.Next = m
	}
	x.last = m

	if x.m == nil {
		x.m = make(map[types.SmiIdentifier]*Macro)
	}
	x.m[m.Name] = m
}

func (x *MacroMap) Get(name types.SmiIdentifier) *Macro {
	if x.m == nil {
		return nil
	}
	return x.m[name]
}

func (x *MacroMap) GetName(name string) *Macro {
	return x.Get(types.SmiIdentifier(name))
}
