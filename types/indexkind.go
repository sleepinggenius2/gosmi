package types

//go:generate enumer -type=IndexKind -autotrimprefix -json

type IndexKind int

const (
	IndexUnknown IndexKind = iota
	IndexIndex
	IndexAugment
	IndexReorder
	IndexSparse
	IndexExpand
)
