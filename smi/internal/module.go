package internal

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/sleepinggenius2/gosmi/parser"
	"github.com/sleepinggenius2/gosmi/types"
)

type Module struct {
	types.SmiModule
	LastUpdated            time.Time
	Identity               *Object
	Objects                ObjectMap
	Types                  TypeMap
	Macros                 MacroMap
	Imports                ImportMap
	FirstRevision          *Revision
	LastRevision           *Revision
	Flags                  Flags
	NumImportedIdentifiers int
	NumStatements          int
	NumModuleIdentities    int
	Prev                   *Module
	Next                   *Module
	PrefixNode             *Node

	pending map[types.SmiIdentifier]*Object
}

func (x *Module) addPending(name types.SmiIdentifier) *Object {
	if x.pending == nil {
		x.pending = make(map[types.SmiIdentifier]*Object)
	}
	obj := new(Object)
	x.pending[name] = obj
	return obj
}

func (x *Module) getPending(name types.SmiIdentifier) *Object {
	if x.pending == nil {
		return nil
	}
	return x.pending[name]
}

func (x *Module) AddRevision(revision *Revision) {
	if x.LastRevision == nil {
		x.FirstRevision = revision
		x.LastRevision = revision
		return
	}
	r := x.LastRevision
	for r != nil && r.Date.After(revision.Date) {
		r = r.Prev
	}
	if r == nil {
		x.FirstRevision.Prev = revision
		revision.Next = x.FirstRevision
		x.FirstRevision = revision
	} else {
		revision.Next = r.Next
		revision.Prev = r
		if r.Next == nil {
			x.LastRevision = revision
		} else {
			r.Next.Prev = revision
		}
		r.Next = revision
	}
}

func (x *Module) GetObject(name types.SmiIdentifier) *Object {
	obj := x.Objects.Get(name)
	if obj != nil {
		return obj
	}
	obj = x.getPending(name)
	if obj != nil {
		return obj
	}
	wellKnown := smiHandle.Modules.Get(WellKnownModuleName)
	if wellKnown != nil {
		obj = wellKnown.Objects.Get(name)
		if obj != nil {
			return obj
		}
	}
	i := x.Imports.Get(name)
	if i == nil {
		return x.addPending(name)
	}
	i.Used = true
	module, err := GetModule(i.Module.String())
	if err != nil {
		return nil
	}
	return module.GetObject(name)
}

func (x *Module) GetType(name types.SmiIdentifier) *Type {
	t := x.Types.Get(name)
	if t != nil {
		return t
	}
	i := x.Imports.Get(name)
	if i == nil {
		return nil
	}
	i.Used = true
	module, err := GetModule(i.Module.String())
	if err != nil {
		return nil
	}
	return module.GetType(i.Name)
}

func (x *Module) IsWellKnown() bool {
	return x != nil && x.Name == WellKnownModuleName
}

func (x *Module) SetPrefixNode(n *Node) {
	if x.PrefixNode == nil {
		x.PrefixNode = n
		return
	}
	if len(n.Oid) < len(x.PrefixNode.Oid) {
		nodePtr := FindNodeByOid(len(n.Oid), x.PrefixNode.Oid)
		if nodePtr == nil {
			// Incomplete object tree
			return
		}
		x.PrefixNode = nodePtr
	}
	for i, subId := range x.PrefixNode.Oid {
		if subId != n.Oid[i] {
			x.PrefixNode = FindNodeByOid(i, x.PrefixNode.Oid)
			return
		}
	}
}

type ModuleMap struct {
	First *Module

	last      *Module
	m         map[types.SmiIdentifier]*Module
	wellKnown *Module
}

func (x *ModuleMap) Add(m *Module) {
	if m.IsWellKnown() {
		x.wellKnown = m
	}
	m.Prev = x.last
	if x.First == nil {
		x.First = m
	} else {
		x.last.Next = m
	}
	x.last = m

	if x.m == nil {
		x.m = make(map[types.SmiIdentifier]*Module)
	}
	x.m[m.Name] = m
}

func (x *ModuleMap) Get(name types.SmiIdentifier) *Module {
	if name == WellKnownModuleName {
		return x.wellKnown
	}
	if x.m == nil {
		return nil
	}
	return x.m[name]
}

func (x *ModuleMap) GetName(name string) *Module {
	return x.Get(types.SmiIdentifier(name))
}

type Import struct {
	types.SmiImport
	ModulePtr *Module
	Flags     Flags
	Prev      *Import
	Next      *Import
	Kind      Kind
	Used      bool
	Line      int
}

type ImportMap struct {
	First *Import

	last *Import
	m    map[types.SmiIdentifier]*Import
}

func (x *ImportMap) Add(i *Import) {
	i.Prev = x.last
	if x.last == nil {
		x.First = i
	} else {
		x.last.Next = i
	}
	x.last = i

	if x.m == nil {
		x.m = make(map[types.SmiIdentifier]*Import)
	}
	x.m[i.Name] = i
	if newImport, ok := importConversions[i.SmiImport]; ok {
		i.Module = newImport.Module
		i.Name = newImport.Name
	}
}

func (x *ImportMap) Get(name types.SmiIdentifier) *Import {
	if x.m == nil {
		return nil
	}
	return x.m[name]
}

func (x *ImportMap) GetName(name string) *Import {
	return x.Get(types.SmiIdentifier(name))
}

type Revision struct {
	types.SmiRevision
	Module *Module
	Prev   *Revision
	Next   *Revision
	Line   int
}

func FindModuleByName(modulename string) *Module {
	return smiHandle.Modules.GetName(modulename)
}

func GetModulePath(name string) (string, error) {
	if name == "" {
		return "", errors.New("Name is required")
	}
	// Relative or absolute
	if name[0] == '.' || name[0] == '~' || filepath.IsAbs(name) {
		dir, file := filepath.Split(name)
		dir, err := expandPath(dir)
		if err != nil {
			return "", errors.Wrap(err, "Expand path")
		}
		return filepath.Join(dir, file), nil
	}

	if filepath.Ext(name) != "" {
		// Filename w/ extension
		for _, path := range smiHandle.Paths {
			fullpath := filepath.Join(path, name)
			info, err := os.Stat(fullpath)
			if err == nil && info.Mode().IsRegular() {
				return fullpath, nil
			}
		}
		return "", os.ErrNotExist
	}

	var modulePath string
	for _, p := range smiHandle.Paths {
		err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if info.Name() == name {
				modulePath = path
				return filepath.SkipDir
			}
			parts := strings.SplitN(info.Name(), ".", 2)
			if parts[0] == name {
				var ext string
				if len(parts) > 1 {
					ext = parts[1]
				}
				if ext == "" || ext == "mib" || ext == "my" || ext == "mi2" || ext == "txt" {
					modulePath = path
					return filepath.SkipDir
				}
			}
			return nil
		})
		if err != nil || modulePath != "" {
			return modulePath, errors.Wrapf(err, "Walk path '%s'", p)
		}
	}
	return "", os.ErrNotExist
}

func GetModule(name string) (*Module, error) {
	module := FindModuleByName(name)
	if module != nil {
		return module, nil
	}
	return LoadModule(name)
}

func LoadModule(name string) (*Module, error) {
	//log.Printf("%s: Loading", name)
	path, err := GetModulePath(name)
	if err != nil {
		return nil, errors.Wrap(err, "Get module path")
	}
	//log.Printf("%s: Found at %s", name, path)
	in, err := parser.ParseFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "Parse module")
	}
	//log.Printf("%s: Parsed", name)
	out, err := BuildModule(path, in)
	if err != nil {
		return nil, errors.Wrap(err, "Build module")
	}
	//log.Printf("%s: Built", name)
	return out, nil
}

type columnMap struct {
	m map[types.SmiIdentifier]struct{}
}

func (x *columnMap) Add(name types.SmiIdentifier) {
	if x.m == nil {
		x.m = make(map[types.SmiIdentifier]struct{})
	}
	x.m[name] = struct{}{}
}

func (x *columnMap) CheckAndRemove(name types.SmiIdentifier) bool {
	_, ok := x.m[name]
	if ok {
		delete(x.m, name)
	}
	return ok
}

func BuildModule(path string, in *parser.Module) (*Module, error) {
	var columnMap columnMap
	out := &Module{
		SmiModule: types.SmiModule{
			Name: in.Name,
			Path: path,
		},
	}

	var currImport *Import
	for _, i := range in.Body.Imports {
		for _, name := range i.Names {
			currImport = &Import{
				SmiImport: types.SmiImport{
					Module: i.Module,
					Name:   name,
				},
				ModulePtr: out,
				Line:      i.Pos.Line,
			}
			out.Imports.Add(currImport)
			out.NumImportedIdentifiers++
		}
	}

	if in.Body.Identity != nil {
		out.NumModuleIdentities = 1
		out.LastUpdated = in.Body.Identity.LastUpdated.ToTime()
		out.Organization = in.Body.Identity.Organization
		out.ContactInfo = in.Body.Identity.ContactInfo
		out.Description = in.Body.Identity.Description
		out.Language = types.LanguageSMIv2

		out.Identity = &Object{
			SmiNode: types.SmiNode{
				Name:        in.Body.Identity.Name,
				Decl:        types.DeclModuleIdentity,
				Description: in.Body.Identity.Description,
				NodeKind:    types.NodeNode,
			},
			Module: out,
		}
		out.Objects.AddWithOid(out.Identity, in.Body.Identity.Oid)

		var currRevision *Revision
		for _, revision := range in.Body.Identity.Revisions {
			currRevision = &Revision{
				SmiRevision: types.SmiRevision{
					Date:        revision.Date.ToTime(),
					Description: revision.Description,
				},
				Module: out,
				Line:   revision.Pos.Line,
			}
			out.AddRevision(currRevision)
		}
	} else {
		out.Language = types.LanguageSMIv1
	}

	var currType *Type
	for _, t := range in.Body.Types {
		if t.Sequence != nil {
			for _, col := range t.Sequence.Entries {
				columnMap.Add(col.Descriptor)
			}
			continue
		}
		currType = &Type{
			SmiType: types.SmiType{
				Name: t.Name,
			},
			Module: out,
			Line:   t.Pos.Line,
		}
		var syntax parser.SyntaxType
		if t.TextualConvention != nil {
			syntax = t.TextualConvention.Syntax
			currType.Decl = types.DeclTextualConvention
			currType.Description = t.TextualConvention.Description
			currType.Format = t.TextualConvention.DisplayHint
			currType.Reference = t.TextualConvention.Reference
			currType.Status = t.TextualConvention.Status.ToSmi()
		} else if t.Implicit != nil {
			syntax = t.Implicit.Syntax
			currType.Decl = types.DeclTypeAssignment
		} else {
			syntax = *t.Syntax
			currType.Decl = types.DeclTypeAssignment
		}
		parentType := GetBaseTypeFromSyntax(syntax)
		if parentType == nil {
			parentType = out.GetType(syntax.Name)
			if parentType == nil {
				// What do we do here?
				break
			}
		}
		if parentType.Decl == types.DeclTextualConvention {
			// This is invalid
			break
		}
		currType.BaseType = parentType.BaseType
		currType.Parent = parentType
		if syntax.SubType != nil && currType.Name != "Integer32" {
			var ranges []parser.Range
			baseType := currType.BaseType
			if baseType == types.BaseTypeOctetString {
				ranges = syntax.SubType.OctetString
				baseType = types.BaseTypeUnsigned32
			} else {
				ranges = syntax.SubType.Integer
				if baseType == types.BaseTypeBits {
					baseType = types.BaseTypeUnsigned32
				}
			}
			rangeSort(ranges)
			for _, r := range ranges {
				if r.End == "" {
					r.End = r.Start
				}
				currType.AddRange(GetValue(r.Start, baseType), GetValue(r.End, baseType))
			}
		} else if len(syntax.Enum) > 0 {
			namedNumberSort(syntax.Enum)
			for _, nn := range syntax.Enum {
				currType.AddNamedNumber(nn.Name, GetValue(nn.Value, currType.BaseType))
			}
			currType.BaseType = types.BaseTypeEnum
		}
		out.Types.Add(currType)
	}

	var currMacro *Macro
	for _, macro := range in.Body.Macros {
		currMacro = &Macro{
			SmiMacro: types.SmiMacro{
				Name: macro.Name,
				Decl: types.DeclMacro,
			},
			Line: macro.Pos.Line,
		}
		out.Macros.Add(currMacro)
	}

	var currObject *Object
	for _, node := range in.Body.Nodes {
		currObject = out.getPending(node.Name)
		if currObject == nil {
			currObject = new(Object)
		}
		currObject.Name = node.Name
		currObject.Module = out
		currObject.Line = node.Pos.Line

		switch {
		case node.ObjectIdentifier:
			currObject.Decl = types.DeclValueAssignment
			currObject.NodeKind = types.NodeNode
		case node.ObjectIdentity != nil:
			currObject.Decl = types.DeclObjectIdentity
			currObject.NodeKind = types.NodeNode
			currObject.Status = node.ObjectIdentity.Status.ToSmi()
			currObject.Description = node.ObjectIdentity.Description
			currObject.Reference = node.ObjectIdentity.Reference
		case node.ObjectGroup != nil:
			currObject.Decl = types.DeclObjectGroup
			currObject.NodeKind = types.NodeGroup
			currObject.Status = node.ObjectGroup.Status.ToSmi()
			currObject.Description = node.ObjectGroup.Description
			currObject.Reference = node.ObjectGroup.Reference
			currObject.AddElements(node.ObjectGroup.Objects)
		case node.ObjectType != nil:
			objType := node.ObjectType
			currObject.Decl = types.DeclObjectType
			currObject.Access = objType.Access.ToSmi()
			currObject.Create = objType.Access == parser.AccessReadCreate
			currObject.Status = objType.Status.ToSmi()
			currObject.Units = objType.Units
			currObject.Description = objType.Description
			currObject.Reference = objType.Reference
			if len(objType.Index) > 0 {
				currObject.NodeKind = types.NodeRow
				currObject.IndexKind = types.IndexIndex
				currObject.Implied = objType.Index[len(objType.Index)-1].Implied
				indices := make([]types.SmiIdentifier, len(objType.Index))
				for i, index := range objType.Index {
					indices[i] = index.Name
				}
				currObject.AddElements(indices)
			} else if objType.Augments != nil {
				currObject.NodeKind = types.NodeRow
				currObject.IndexKind = types.IndexAugment
				currObject.Related = out.GetObject(*objType.Augments)
			} else if objType.Syntax.Sequence != nil {
				currObject.NodeKind = types.NodeTable
			} else {
				if columnMap.CheckAndRemove(node.Name) {
					currObject.NodeKind = types.NodeColumn
				} else {
					currObject.NodeKind = types.NodeScalar
				}
				syntax := *objType.Syntax.Type
				parentType := GetBaseTypeFromSyntax(syntax)
				if parentType == nil {
					parentType = out.GetType(syntax.Name)
					if parentType == nil {
						// What do we do here?
						break
					}
				}
				if syntax.SubType == nil && len(syntax.Enum) == 0 {
					currObject.Type = parentType
					break
				}
				currType = &Type{
					SmiType: types.SmiType{
						BaseType: parentType.BaseType,
						Decl:     types.DeclImplicitType,
						Status:   currObject.Status,
					},
					Module: out,
					Parent: parentType,
					Line:   syntax.Pos.Line,
				}
				baseType := currType.BaseType
				if syntax.SubType != nil {
					var ranges []parser.Range
					if baseType == types.BaseTypeOctetString {
						ranges = syntax.SubType.OctetString
						baseType = types.BaseTypeUnsigned32
					} else {
						ranges = syntax.SubType.Integer
					}
					rangeSort(ranges)
					for _, r := range ranges {
						if r.End == "" {
							r.End = r.Start
						}
						currType.AddRange(GetValue(r.Start, baseType), GetValue(r.End, baseType))
					}
				} else if len(syntax.Enum) > 0 {
					if baseType == types.BaseTypeEnum {
						if parentType.List == nil || parentType.List.Ptr == nil {
							// TODO: Figure out a better option. This should never happen.
							baseType = types.BaseTypeInteger32
						} else {
							baseType = parentType.List.Ptr.(*NamedNumber).Value.BaseType
						}
					} else if baseType == types.BaseTypeBits {
						baseType = types.BaseTypeUnsigned32
					}
					namedNumberSort(syntax.Enum)
					for _, nn := range syntax.Enum {
						currType.AddNamedNumber(nn.Name, GetValue(nn.Value, baseType))
					}
					if currType.BaseType == types.BaseTypeBits {
						if parentType == smiHandle.TypeBits {
							currType.Name = "Bits"
						} else {
							currType.Name = parentType.Name
						}
					} else {
						if parentType.Module == nil || parentType.Module.IsWellKnown() {
							currType.Name = "Enumeration"
						} else {
							currType.Name = parentType.Name
						}
						currType.BaseType = types.BaseTypeEnum
					}
				}
				currObject.Type = currType
			}
		case node.NotificationGroup != nil:
			currObject.Decl = types.DeclNotificationGroup
			currObject.NodeKind = types.NodeGroup
			currObject.Status = node.NotificationGroup.Status.ToSmi()
			currObject.Description = node.NotificationGroup.Description
			currObject.Reference = node.NotificationGroup.Reference
			currObject.AddElements(node.NotificationGroup.Notifications)
		case node.NotificationType != nil:
			currObject.Decl = types.DeclNotificationType
			currObject.NodeKind = types.NodeNotification
			currObject.Status = node.NotificationType.Status.ToSmi()
			currObject.Description = node.NotificationType.Description
			currObject.Reference = node.NotificationType.Reference
			currObject.AddElements(node.NotificationType.Objects)
		case node.ModuleCompliance != nil:
			currObject.Decl = types.DeclModuleCompliance
			currObject.NodeKind = types.NodeCompliance
			currObject.Status = node.ModuleCompliance.Status.ToSmi()
			currObject.Description = node.ModuleCompliance.Description
			currObject.Reference = node.ModuleCompliance.Reference
			// TODO: Deal with node.ModuleCompliance.Modules
		case node.AgentCapabilities != nil:
			currObject.Decl = types.DeclAgentCapabilities
			currObject.NodeKind = types.NodeCapabilities
			currObject.Status = node.AgentCapabilities.Status.ToSmi()
			currObject.Description = node.AgentCapabilities.Description
			currObject.Reference = node.AgentCapabilities.Reference
			// TODO: Deal with node.AgentCapabilities.Modules
		case node.TrapType != nil:
			currObject.Decl = types.DeclTrapType
			currObject.NodeKind = types.NodeNotification
			currObject.Description = node.TrapType.Description
			currObject.Reference = node.TrapType.Reference
			currObject.AddElements(node.TrapType.Objects)
			node.Oid = &parser.Oid{
				SubIdentifiers: []parser.SubIdentifier{
					parser.SubIdentifier{Name: &node.TrapType.Enterprise},
					parser.SubIdentifier{Number: new(types.SmiSubId)},
					parser.SubIdentifier{Number: node.SubIdentifier},
				},
			}
		default:
			// This should never happen
			return nil, errors.Errorf("Cannot determine type for node '%s'", node.Name)
		}
		out.Objects.AddWithOid(currObject, *node.Oid)
	}
	smiHandle.Modules.Add(out)
	return out, nil
}
