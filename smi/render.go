package smi

import (
	"strconv"
	"strings"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

func RenderNode(smiNodePtr *types.SmiNode, flags types.Render) string {
	if smiNodePtr == nil {
		if flags&types.RenderUnknown > 0 {
			return internal.UnknownLabel
		}
		return ""
	}
	modulePtr := GetNodeModule(smiNodePtr)
	if flags&types.RenderQualified == 0 || modulePtr == nil || modulePtr.Name == "" {
		return smiNodePtr.Name.String()
	}
	return modulePtr.Name.String() + "::" + smiNodePtr.Name.String()
}

func RenderOID(oid types.Oid, flags types.Render) string {
	if len(oid) == 0 {
		if flags&types.RenderUnknown > 0 {
			return internal.UnknownLabel
		}
		return ""
	}
	var i int
	var b strings.Builder
	if flags&(types.RenderName|types.RenderQualified) > 0 {
		nodePtr := GetNodeByOID(oid)
		if nodePtr != nil && nodePtr.Name != "" {
			i = nodePtr.OidLen
			b.WriteString(RenderNode(nodePtr, flags))
		}
	}
	for ; i < len(oid); i++ {
		if b.Len() > 0 {
			b.WriteRune('.')
		}
		b.WriteString(strconv.FormatUint(uint64(oid[i]), 10))
	}
	return b.String()
}
