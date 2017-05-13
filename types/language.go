package types

//go:generate enumer -type=Language -autotrimprefix -json

type Language int

const (
	LanguageUnknown Language = iota
	LanguageSMIv1
	LanguageSMIv2
	LanguageSMIng
	LanguageSPPI
)
