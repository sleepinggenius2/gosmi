package models

import (
	"github.com/sleepinggenius2/gosmi"
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

func (t TableNode) Index() []ColumnNode {
	return t.Row.Index
}

type RowNode struct {
	BaseNode
	Columns []ColumnNode
	Index   []ColumnNode
}

type ColumnNode struct {
	ScalarNode
}

type NotificationNode struct {
	BaseNode
	Objects []ScalarNode
}
