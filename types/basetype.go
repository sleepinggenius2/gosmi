package types

//go:generate enumer -type=BaseType -autotrimprefix -json

type BaseType int

const (
	BaseTypeUnknown BaseType = iota
	BaseTypeInteger32
	BaseTypeOctetString
	BaseTypeObjectIdentifier
	BaseTypeUnsigned32
	BaseTypeInteger64
	BaseTypeUnsigned64
	BaseTypeFloat32
	BaseTypeFloat64
	BaseTypeFloat128
	BaseTypeEnum
	BaseTypeBits
	BaseTypePointer
)
