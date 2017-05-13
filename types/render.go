package types

//go:generate enumer -type=Render -autotrimprefix -json

type Render int

const (
	RenderNumeric Render = 1 << iota
	RenderName
	RenderQualified
	RenderFormat
	RenderPrintable
	RenderUnknown
	RenderAll Render = 0xff
)
