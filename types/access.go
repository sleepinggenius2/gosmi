package types

//go:generate enumer -type=Access -autotrimprefix -json

type Access int

const (
	AccessUnknown Access = iota
	AccessNotImplemented
	AccessNotAccessible
	AccessNotify
	AccessReadOnly
	AccessReadWrite
	AccessInstall
	AccessInstallNotify
	AccessReportOnly
	AccessEventOnly
)
