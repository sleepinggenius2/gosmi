package models

import (
	"github.com/pkg/errors"

	"github.com/sleepinggenius2/gosmi/types"
)

type BaseNode struct {
	Name         string
	Oid          types.Oid
	OidFormatted string
	OidLen       uint
}

func (b BaseNode) ChildOf(n BaseNode) bool {
	return b.Oid.ChildOf(n.Oid)
}

func (b BaseNode) ParentOf(n BaseNode) bool {
	return b.Oid.ParentOf(n.Oid)
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

func (t TableNode) Implied() bool {
	return t.Row.Implied
}

func (t TableNode) Index() []ColumnNode {
	return t.Row.Index
}

func (t TableNode) BuildIndex(index ...interface{}) (types.Oid, error) {
	tableIndex := t.Row.Index
	indexLen := len(index)
	if indexLen == 0 {
		return nil, nil
	}
	tableIndexLen := len(tableIndex)
	if indexLen > tableIndexLen {
		return nil, errors.New("Too many index values given")
	}
	if v, ok := index[0].(types.Oid); ok {
		return v, nil
	}
	ret := make(types.Oid, 0, len(index))
	for i := range index {
		indexValue, err := tableIndex[i].Type.IndexValue(index[i], t.Row.Implied && (i == tableIndexLen-1))
		if err != nil {
			return nil, errors.Wrap(err, tableIndex[i].Name+": "+tableIndex[i].Type.BaseType.String())
		}
		ret = append(ret, indexValue...)
	}
	return ret, nil
}

type RowNode struct {
	BaseNode
	Columns []ColumnNode
	Implied bool
	Index   []ColumnNode
}

type ColumnNode ScalarNode

type NotificationNode struct {
	BaseNode
	Objects []ScalarNode
}
