package gosmi

import (
	"github.com/sleepinggenius2/gosmi/smi"
	"github.com/sleepinggenius2/gosmi/types"
)

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

func (t SmiNode) getRow() (row *types.SmiNode) {
	switch t.Kind {
	case types.NodeRow:
		row = t.GetRaw()
	case types.NodeTable:
		row = smi.GetFirstChildNode(t.smiNode)
		if row == nil {
			return
		}
	default:
		return
	}

	if row.NodeKind != types.NodeRow {
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

	for smiColumn := smi.GetFirstChildNode(row); smiColumn != nil; smiColumn = smi.GetNextChildNode(smiColumn) {
		if smiColumn.NodeKind != types.NodeColumn {
			// TODO: error
			return
		}
		column := CreateNode(smiColumn)
		columns[column.Name] = column
		columnOrder = append(columnOrder, column.Name)
	}
	return
}

func (t SmiNode) GetImplied() bool {
	row := t.getRow()
	if row == nil {
		return false
	}

	return row.Implied
}

func (t SmiNode) GetAugment() (row SmiNode) {
	smiRow := t.getRow()
	if smiRow == nil {
		return
	}

	if smiRow.IndexKind != types.IndexAugment {
		return
	}

	smiRow = smi.GetRelatedNode(smiRow)
	if smiRow == nil {
		return
	}

	if smiRow.NodeKind != types.NodeRow {
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

	if row.IndexKind == types.IndexAugment {
		row = smi.GetRelatedNode(row)
		if row == nil {
			return
		}

		if row.NodeKind != types.NodeRow {
			// TODO: error
			return
		}
	} else if row.IndexKind != types.IndexIndex {
		// TODO: unsupported
		return
	}

	for smiElement := smi.GetFirstElement(row); smiElement != nil; smiElement = smi.GetNextElement(smiElement) {
		smiColumn := smi.GetElementNode(smiElement)
		if smiColumn == nil {
			// TODO: error
			return
		}
		if smiColumn.NodeKind != types.NodeColumn {
			// TODO: error
			return
		}
		index = append(index, CreateNode(smiColumn))
	}
	return
}
