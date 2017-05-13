package types

//go:generate enumer -type=Status -autotrimprefix -json

type Status int

const (
	StatusUnknown Status = iota
	StatusCurrent
	StatusDeprecated
	StatusMandatory
	StatusOptional
	StatusObsolete
)
