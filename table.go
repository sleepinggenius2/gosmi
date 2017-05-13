package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

import "github.com/sleepinggenius2/gosmi/types"

func (t Node) getRow() (row *C.struct_SmiNode) {
	row = C.smiGetFirstChildNode(t.SmiNode)
	if row == nil {
		return
	}

	if types.NodeKind(row.nodekind) != types.NodeRow {
		// TODO: error
		return nil
	}

	return
}

func (t Node) GetColumns() (columns []Node) {
	row := t.getRow()
	if row == nil {
		return
	}

	for column := C.smiGetFirstChildNode(row); column != nil; column = C.smiGetNextChildNode(column) {
		if types.NodeKind(column.nodekind) != types.NodeColumn {
			// TODO: error
			return
		}
		columns = append(columns, CreateNode(column))
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

	for element := C.smiGetFirstElement(row); element != nil; element = C.smiGetNextElement(element) {
		column := C.smiGetElementNode(element)
		if column == nil {
			// TODO: error
			return
		}
		if types.NodeKind(column.nodekind) != types.NodeColumn {
			// TODO: error
			return
		}
		index = append(index, CreateNode(column))
	}
	return
}
