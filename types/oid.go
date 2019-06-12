package types

import (
	"strconv"
	"strings"
)

type SmiSubId uint32

type Oid []SmiSubId

func (o Oid) After(oid Oid) bool {
	myLen := len(o)
	oidLen := len(oid)
	for i := 0; i < oidLen; i++ {
		if myLen == i {
			return false
		}
		if o[i] != oid[i] {
			return o[i] > oid[i]
		}
	}
	return myLen > oidLen
}

func (o Oid) Before(oid Oid) bool {
	return oid.After(o)
}

func (o Oid) ChildOf(oid Oid) bool {
	myLen := len(o)
	oidLen := len(oid)
	for i := 0; i < oidLen; i++ {
		if myLen == i {
			return false
		}
		if o[i] != oid[i] {
			return false
		}
	}
	return true
}

func (o Oid) Equals(oid Oid) bool {
	if len(o) != len(oid) {
		return false
	}
	for i := range o {
		if o[i] != oid[i] {
			return false
		}
	}
	return true
}

func (o Oid) ParentOf(oid Oid) bool {
	return oid.ChildOf(o)
}

func (o Oid) String() string {
	oidParts := make([]string, len(o))
	for i, oidPart := range o {
		oidParts[i] = strconv.FormatUint(uint64(oidPart), 10)
	}
	return strings.Join(oidParts, ".")
}

func NewOid(parent Oid, subId SmiSubId) Oid {
	oid := make(Oid, len(parent), len(parent)+1)
	copy(oid, parent)
	return append(oid, subId)
}

func OidFromString(s string) (Oid, error) {
	oidParts := strings.Split(strings.Trim(strings.TrimSpace(s), "."), ".")
	oid := make(Oid, len(oidParts))
	for i := range oidParts {
		oidPart, err := strconv.ParseUint(oidParts[i], 10, 32)
		if err != nil {
			return nil, err
		}
		oid[i] = SmiSubId(oidPart)
	}
	return oid, nil
}

// Helper for defining constants from strings
func OidMustFromString(s string) Oid {
	oid, err := OidFromString(s)
	if err != nil {
		panic(err)
	}
	return oid
}
