package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import "github.com/sleepinggenius2/gosmi/types"

type Table struct {
	Node
	Columns map[string]Node
	Implied bool
	Index []Node
}

func (t Node) AsTable() Table {
	return Table{
		Node: t,
		Columns: t.GetColumns(),
		Implied: t.GetImplied(),
		Index: t.GetIndex(),
	}
}

func (t Node) getRow() (row *C.struct_SmiNode) {
	row = C.smiGetFirstChildNode(t.smiNode)
	if row == nil {
		return
	}

	if types.NodeKind(row.nodekind) != types.NodeRow {
		// TODO: error
		return nil
	}

	return
}

func (t Node) GetColumns() (columns map[string]Node) {
	row := t.getRow()
	if row == nil {
		return
	}

	columns = make(map[string]Node)

	for smiColumn := C.smiGetFirstChildNode(row); smiColumn != nil; smiColumn = C.smiGetNextChildNode(smiColumn) {
		if types.NodeKind(smiColumn.nodekind) != types.NodeColumn {
			// TODO: error
			return
		}
		column := CreateNode(smiColumn)
		columns[column.Name] = column
	}
	return
}

func (t Node) GetImplied() (implied bool) {
	row := t.getRow()
	if row == nil {
		return false
	}

	return int(row.implied) > 0
}

func (t Node) GetIndex() (index []Node) {
	row := t.getRow()
	if row == nil {
		return
	}

	if types.IndexKind(row.indexkind) != types.IndexIndex {
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
