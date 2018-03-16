package models

type BaseNode struct {
	Name         string
	Oid          []uint
	OidFormatted string
	OidLen       uint
}

type ScalarNode struct {
	BaseNode
	Type Type
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

type ColumnNode ScalarNode

type NotificationNode struct {
	BaseNode
	Objects []ScalarNode
}
