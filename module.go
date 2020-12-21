package gosmi

import (
	"fmt"

	"github.com/sleepinggenius2/gosmi/models"
	"github.com/sleepinggenius2/gosmi/smi"
	"github.com/sleepinggenius2/gosmi/types"
)

type SmiModule struct {
	models.Module
	smiModule *types.SmiModule
}

func (m SmiModule) GetIdentityNode() (node SmiNode, ok bool) {
	smiIdentityNode := smi.GetModuleIdentityNode(m.smiModule)
	if smiIdentityNode == nil {
		return
	}
	return CreateNode(smiIdentityNode), true
}

func (m SmiModule) GetImports() (imports []models.Import) {
	for smiImport := smi.GetFirstImport(m.smiModule); smiImport != nil; smiImport = smi.GetNextImport(smiImport) {
		_import := models.Import{
			Module: string(smiImport.Module),
			Name:   string(smiImport.Name),
		}
		imports = append(imports, _import)
	}
	return
}

func (m SmiModule) GetNode(name string) (node SmiNode, err error) {
	return GetNode(name, m)
}

func (m SmiModule) GetNodes(kind ...types.NodeKind) (nodes []SmiNode) {
	nodeKind := types.NodeAny
	if len(kind) > 0 && kind[0] != types.NodeUnknown {
		nodeKind = kind[0]
	}
	for smiNode := smi.GetFirstNode(m.smiModule, nodeKind); smiNode != nil; smiNode = smi.GetNextNode(smiNode, nodeKind) {
		nodes = append(nodes, CreateNode(smiNode))
	}
	return
}

func (m SmiModule) GetRevisions() (revisions []models.Revision) {
	for smiRevision := smi.GetFirstRevision(m.smiModule); smiRevision != nil; smiRevision = smi.GetNextRevision(smiRevision) {
		revision := models.Revision{
			Date:        smiRevision.Date,
			Description: smiRevision.Description,
		}
		revisions = append(revisions, revision)
	}
	return
}

func (m SmiModule) GetType(name string) (outType SmiType, err error) {
	return GetType(name, m)
}

func (m SmiModule) GetTypes() (types []SmiType) {
	for smiType := smi.GetFirstType(m.smiModule); smiType != nil; smiType = smi.GetNextType(smiType) {
		types = append(types, CreateType(smiType))
	}
	return
}

func (m SmiModule) GetRaw() (module *types.SmiModule) {
	return m.smiModule
}

func (m *SmiModule) SetRaw(smiModule *types.SmiModule) {
	m.smiModule = smiModule
}

func CreateModule(smiModule *types.SmiModule) (module SmiModule) {
	return SmiModule{
		Module: models.Module{
			ContactInfo:  smiModule.ContactInfo,
			Description:  smiModule.Description,
			Language:     smiModule.Language,
			Name:         string(smiModule.Name),
			Organization: smiModule.Organization,
			Path:         smiModule.Path,
			Reference:    smiModule.Reference,
		},
		smiModule: smiModule,
	}
}

func LoadModule(modulePath string) (string, error) {
	moduleName := smi.LoadModule(modulePath)
	if moduleName == "" {
		return "", fmt.Errorf("Could not load module at %s", modulePath)
	}
	return moduleName, nil
}

func GetLoadedModules() (modules []SmiModule) {
	for smiModule := smi.GetFirstModule(); smiModule != nil; smiModule = smi.GetNextModule(smiModule) {
		modules = append(modules, CreateModule(smiModule))
	}
	return
}

func IsLoaded(moduleName string) bool {
	return smi.IsLoaded(moduleName)
}

func GetModule(name string) (module SmiModule, err error) {
	smiModule := smi.GetModule(name)
	if smiModule == nil {
		err = fmt.Errorf("Could not find module named %s", name)
		return
	}
	return CreateModule(smiModule), nil
}
