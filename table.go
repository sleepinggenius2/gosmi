package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import "github.com/sleepinggenius2/gosmi/types"

type Table struct {
	SmiNode
	Columns     map[string]SmiNode
	ColumnOrder []string
	Implied     bool
	Index       []SmiNode
}

func (t SmiNode) AsTable() Table {
	columns, columnOrder := t.GetColumns()
	return Table{
		SmiNode:     t,
		Columns:     columns,
		ColumnOrder: columnOrder,
		Implied:     t.GetImplied(),
		Index:       t.GetIndex(),
	}
}

func (t SmiNode) getRow() (row *C.struct_SmiNode) {
	switch t.Kind {
	case types.NodeRow:
		row = t.GetRaw()
	case types.NodeTable:
		row = C.smiGetFirstChildNode(t.smiNode)
		if row == nil {
			return
		}
	default:
		return
	}

	if types.NodeKind(row.nodekind) != types.NodeRow {
		// TODO: error
		return nil
	}

	return
}

func (t SmiNode) GetRow() (row SmiNode) {
	smiRow := t.getRow()
	if smiRow == nil {
		return
	}
	return CreateNode(smiRow)
}

func (t SmiNode) GetColumns() (columns map[string]SmiNode, columnOrder []string) {
	row := t.getRow()
	if row == nil {
		return
	}

	columns = make(map[string]SmiNode)
	columnOrder = make([]string, 0, 2)

	for smiColumn := C.smiGetFirstChildNode(row); smiColumn != nil; smiColumn = C.smiGetNextChildNode(smiColumn) {
		if types.NodeKind(smiColumn.nodekind) != types.NodeColumn {
			// TODO: error
			return
		}
		column := CreateNode(smiColumn)
		columns[column.Name] = column
		columnOrder = append(columnOrder, column.Name)
	}
	return
}

func (t SmiNode) GetImplied() (implied bool) {
	row := t.getRow()
	if row == nil {
		return false
	}

	return int(row.implied) > 0
}

func (t SmiNode) GetAugment() (row SmiNode) {
	smiRow := t.getRow()
	if smiRow == nil {
		return
	}

	if types.IndexKind(smiRow.indexkind) != types.IndexAugment {
		return
	}

	smiRow = C.smiGetRelatedNode(smiRow)
	if smiRow == nil {
		return
	}

	if types.NodeKind(smiRow.nodekind) != types.NodeRow {
		// TODO: error
		return
	}

	return CreateNode(smiRow)
}

func (t SmiNode) GetIndex() (index []SmiNode) {
	row := t.getRow()
	if row == nil {
		return
	}

	if types.IndexKind(row.indexkind) == types.IndexAugment {
		row = C.smiGetRelatedNode(row)
		if row == nil {
			return
		}

		if types.NodeKind(row.nodekind) != types.NodeRow {
			// TODO: error
			return
		}
	} else if types.IndexKind(row.indexkind) != types.IndexIndex {
		// TODO: unsupported
		return
	}

	for smiElement := C.smiGetFirstElement(row); smiElement != nil; smiElement = C.smiGetNextElement(smiElement) {
		smiColumn := C.smiGetElementNode(smiElement)
		if smiColumn == nil {
			// TODO: error
			return
		}
		if types.NodeKind(smiColumn.nodekind) != types.NodeColumn {
			// TODO: error
			return
		}
		index = append(index, CreateNode(smiColumn))
	}
	return
}
