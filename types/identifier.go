package types

type SmiIdentifier string

func (x SmiIdentifier) String() string {
	if x == "OBJECT IDENTIFIER" {
		return "ObjectIdentifier"
	} else if x == "OCTET STRING" {
		return "OctetString"
	}
	return string(x)
}
