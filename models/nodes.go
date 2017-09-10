package models

import (
	"github.com/sleepinggenius2/gosmi"
	"github.com/sleepinggenius2/gosmi/types"
)

type BaseNode struct {
	Name   string
	Oid    []uint
	OidLen uint
}

type ScalarNode struct {
	BaseNode
	Type gosmi.Type
}

type TableNode struct {
	BaseNode
	Row RowNode
}

func (t TableNode) Columns() []ColumnNode {
	return t.Row.Columns
}

func (t TableNode) Index() []ScalarNode {
	return t.Row.Index
}

type RowNode struct {
	BaseNode
	Columns []ColumnNode
	Index   []ScalarNode
}

type ColumnNode struct {
	ScalarNode
}

type NotificationNode struct {
	BaseNode
	Objects []ScalarNode
}

