package gosmi

import (
	"fmt"

	"github.com/sleepinggenius2/gosmi/models"
	"github.com/sleepinggenius2/gosmi/smi"
	"github.com/sleepinggenius2/gosmi/types"
)

type SmiType struct {
	models.Type
	smiType *types.SmiType
}

func (t *SmiType) getEnum() {
	if t.BaseType == types.BaseTypeUnknown || !(t.BaseType == types.BaseTypeEnum || t.BaseType == types.BaseTypeBits) {
		return
	}

	smiNamedNumber := smi.GetFirstNamedNumber(t.smiType)
	if smiNamedNumber == nil {
		return
	}

	enum := models.Enum{
		BaseType: types.BaseType(smiNamedNumber.Value.BaseType),
	}
	for ; smiNamedNumber != nil; smiNamedNumber = smi.GetNextNamedNumber(smiNamedNumber) {
		namedNumber := models.NamedNumber{
			Name:  string(smiNamedNumber.Name),
			Value: convertValue(smiNamedNumber.Value),
		}
		enum.Values = append(enum.Values, namedNumber)
	}
	t.Enum = &enum
	return
}

func (t SmiType) GetModule() (module SmiModule) {
	smiModule := smi.GetTypeModule(t.smiType)
	return CreateModule(smiModule)
}

func (t *SmiType) getRanges() {
	if t.BaseType == types.BaseTypeUnknown {
		return
	}

	ranges := make([]models.Range, 0)
	// Workaround for libsmi bug that causes ranges to loop infinitely sometimes
	var currSmiRange *types.SmiRange
	for smiRange := smi.GetFirstRange(t.smiType); smiRange != nil && smiRange != currSmiRange; smiRange = smi.GetNextRange(smiRange) {
		r := models.Range{
			BaseType: smiRange.MinValue.BaseType,
			MinValue: convertValue(smiRange.MinValue),
			MaxValue: convertValue(smiRange.MaxValue),
		}
		ranges = append(ranges, r)
		currSmiRange = smiRange
	}
	t.Ranges = ranges
}

func (t SmiType) String() string {
	return t.Type.String()
}

func (t SmiType) GetRaw() (outType *types.SmiType) {
	return t.smiType
}

func (t *SmiType) SetRaw(smiType *types.SmiType) {
	t.smiType = smiType
}

func CreateType(smiType *types.SmiType) (outType SmiType) {
	if smiType == nil {
		return
	}

	outType.SetRaw(smiType)
	outType.BaseType = smiType.BaseType

	if smiType.Name == "" {
		smiType = smi.GetParentType(smiType)
	}

	outType.Decl = smiType.Decl
	outType.Description = smiType.Description
	outType.Format = smiType.Format
	outType.Name = string(smiType.Name)
	outType.Reference = smiType.Reference
	outType.Status = smiType.Status
	outType.Units = smiType.Units

	outType.getEnum()
	outType.getRanges()

	return
}

func CreateTypeFromNode(smiNode *types.SmiNode) (outType *SmiType) {
	smiType := smi.GetNodeType(smiNode)

	if smiType == nil {
		return
	}

	tempType := CreateType(smiType)
	outType = &tempType

	if smiNode.Format != "" {
		outType.Format = smiNode.Format
	}
	if smiNode.Units != "" {
		outType.Units = smiNode.Units
	}

	return
}

func GetType(name string, module ...SmiModule) (outType SmiType, err error) {
	var smiModule *types.SmiModule
	if len(module) > 0 {
		smiModule = module[0].GetRaw()
	}

	smiType := smi.GetType(smiModule, name)
	if smiType == nil {
		if len(module) > 0 {
			err = fmt.Errorf("Could not find type named %s in module %s", name, module[0].Name)
		} else {
			err = fmt.Errorf("Could not find type named %s", name)
		}
		return
	}
	return CreateType(smiType), nil
}

func convertValue(value types.SmiValue) (outValue int64) {
	switch v := value.Value.(type) {
	case int32:
		outValue = int64(v)
	case int64:
		outValue = int64(v)
	case uint32:
		outValue = int64(v)
	case uint64:
		outValue = int64(v)
	}
	return
}
