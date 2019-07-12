package internal

import (
	"github.com/sleepinggenius2/gosmi/parser"
	"github.com/sleepinggenius2/gosmi/types"
)

type Option struct {
	types.SmiOption
	Compliance *Object
	Object     *Object
	Line       int
	List       *List
}

type Refinement struct {
	types.SmiRefinement
	Compliance *Object
	Object     *Object
	Type       *Type
	WriteType  *Type
	Line       int
	List       *List
}

type Index struct {
	Implied   int
	IndexKind types.IndexKind
	List      *List
	Row       *Object
}

type Object struct {
	types.SmiNode
	Module         *Module
	Flags          Flags
	Type           *Type
	Related        *Object
	List           *List
	OptionList     *List
	RefinementList *List
	Node           *Node
	Prev           *Object
	Next           *Object
	PrevSameNode   *Object
	NextSameNode   *Object
	UniquenessPtr  *List
	Line           int

	lastList           *List
	lastOptionList     *List
	lastRefinementList *List
}

func (x *Object) AddElement(obj *Object) {
	list := &List{Ptr: obj}
	if x.lastList == nil {
		x.List = list
	} else {
		x.lastList.Next = list
	}
	x.lastList = list
}

func (x *Object) AddElements(ids []types.SmiIdentifier) {
	for _, objName := range ids {
		obj := x.Module.GetObject(objName)
		// What do we do if it doesn't exist?
		if obj != nil {
			x.AddElement(obj)
		}
	}
}

func (x *Object) AddOption(option *Option) {
	list := &List{Ptr: option}
	option.List = list
	if x.lastOptionList == nil {
		x.OptionList = list
	} else {
		x.lastOptionList.Next = list
	}
	x.lastOptionList = list
}

func (x *Object) AddRefinement(refinement *Refinement) {
	list := &List{Ptr: refinement}
	refinement.List = list
	if x.lastRefinementList == nil {
		x.RefinementList = list
	} else {
		x.lastRefinementList.Next = list
	}
	x.lastRefinementList = list
}

func (x *Object) GetSmiNode() *types.SmiNode {
	if len(x.Oid) == 0 && x.Node != nil {
		x.Oid = x.Node.Oid
	}
	if x.OidLen <= 0 {
		x.OidLen = len(x.Oid)
	}
	return &x.SmiNode
}

type ObjectMap struct {
	First *Object

	last    *Object
	m       map[types.SmiIdentifier]*Object
	pending map[types.SmiIdentifier][]*Node
}

func (x *ObjectMap) linkPending(o *Object) {
	if o.Node == nil {
		return
	}
	if o.Module != nil {
		o.Module.SetPrefixNode(o.Node)
	}
	if x.pending == nil {
		return
	}
	pending := x.pending[o.Name]
	if pending == nil {
		return
	}
	for _, child := range pending {
		child.Parent = o.Node
		o.Node.Children.Add(child)
		for objPtr := child.FirstObject; objPtr != nil; objPtr = objPtr.NextSameNode {
			if objPtr.Module == o.Module {
				x.linkPending(objPtr)
			}
		}
	}
	delete(x.pending, o.Name)
}

func (x *ObjectMap) addPending(name types.SmiIdentifier, node *Node) {
	if x.pending == nil {
		x.pending = make(map[types.SmiIdentifier][]*Node)
	} else {
		for _, n := range x.pending[name] {
			if n.SubId != node.SubId {
				continue
			}
			for objPtr := node.FirstObject; objPtr != nil; objPtr = objPtr.NextSameNode {
				n.AddObject(objPtr)
			}
			return
		}
	}
	x.pending[name] = append(x.pending[name], node)
}

func (x *ObjectMap) linkObject(oid parser.Oid, o *Object) {
	if o.Module == nil || len(oid.SubIdentifiers) == 0 {
		return
	}

	var parentNodePtr *Node
	if o.Module.IsWellKnown() {
		parentNodePtr = smiHandle.RootNode
	} else if oid.SubIdentifiers[0].Name == nil {
		if oid.SubIdentifiers[0].Number == nil {
			return
		}
		wellKnownModule := smiHandle.Modules.Get(WellKnownModuleName)
		if wellKnownModule == nil {
			return
		}
		parentNodePtr = wellKnownModule.PrefixNode.Children.Get(*oid.SubIdentifiers[0].Number)
		if parentNodePtr == nil {
			return
		}
		oid.SubIdentifiers = oid.SubIdentifiers[1:]
	}

	// Final sub-identifier must be a number
	lastSubId := oid.SubIdentifiers[len(oid.SubIdentifiers)-1]
	if lastSubId.Number == nil || lastSubId.Name != nil {
		return
	}

	var nodePtr *Node
	var parentName types.SmiIdentifier
	for i := 0; i < len(oid.SubIdentifiers)-1; i++ {
		subId := oid.SubIdentifiers[i]
		if subId.Name != nil {
			obj := o.Module.GetObject(*subId.Name)
			if obj != nil && obj.Node != nil {
				parentName = ""
				parentNodePtr = obj.Node
				continue
			}
			if subId.Number == nil {
				parentName = *subId.Name
				continue
			}
		}
		// This shouldn't be possible from the parser
		if subId.Number == nil {
			return
		}
		nodePtr = &Node{
			SubId:  *subId.Number,
			Parent: parentNodePtr,
		}
		if subId.Name != nil {
			// Create parent
			parent := &Object{
				SmiNode: types.SmiNode{
					Name:     *subId.Name,
					Decl:     types.DeclImplObject,
					NodeKind: types.NodeNode,
				},
				Module: o.Module,
				Line:   subId.Pos.Line,
				Node:   nodePtr,
			}
			nodePtr.AddObject(parent)
		}
		if parentNodePtr == nil {
			x.addPending(parentName, nodePtr)
			continue
		}
		parentNodePtr.Children.Add(nodePtr)
		parentNodePtr = nodePtr
	}
	o.Node = &Node{
		SubId:       *lastSubId.Number,
		FirstObject: o,
		LastObject:  o,
		Parent:      parentNodePtr,
	}
	if parentNodePtr == nil {
		x.addPending(parentName, o.Node)
	} else {
		parentNodePtr.Children.Add(o.Node)
		x.linkPending(o)
	}
}

func (x *ObjectMap) Add(o *Object) {
	o.Prev = x.last
	if x.First == nil {
		x.First = o
	} else {
		x.last.Next = o
	}
	x.last = o

	if x.m == nil {
		x.m = make(map[types.SmiIdentifier]*Object)
	}
	x.m[o.Name] = o
}

func (x *ObjectMap) AddWithOid(o *Object, oid parser.Oid) {
	x.linkObject(oid, o)
	x.Add(o)
}

func (x *ObjectMap) Get(name types.SmiIdentifier) *Object {
	if x.m == nil {
		return nil
	}
	return x.m[name]
}

func (x *ObjectMap) GetName(name string) *Object {
	return x.Get(types.SmiIdentifier(name))
}

func FindObjectByNode(nodePtr *Node) *Object {
	return nodePtr.FirstObject
}

func FindObjectByModuleAndNode(modulePtr *Module, nodePtr *Node) *Object {
	if modulePtr == nil {
		return nil
	}
	for obj := nodePtr.FirstObject; obj != nil; obj = obj.NextSameNode {
		if obj.Module == modulePtr {
			return obj
		}
	}
	return nil
}

func FindObjectByModuleNameAndNode(module string, nodePtr *Node) *Object {
	return FindObjectByModuleAndNode(FindModuleByName(module), nodePtr)
}

func GetNextChildObject(startNodePtr *Node, modulePtr *Module, nodekind types.NodeKind) *Object {
	if startNodePtr == nil {
		return nil
	}
	var objPtr *Object
	for nodePtr := startNodePtr; nodePtr != nil; nodePtr = nodePtr.Next {
		for objPtr = nodePtr.FirstObject; objPtr != nil; objPtr = objPtr.NextSameNode {
			if (modulePtr == nil || objPtr.Module == modulePtr) && (nodekind == types.NodeAny || (objPtr.NodeKind&nodekind) > 0) {
				return objPtr
			}
		}
		objPtr = GetNextChildObject(nodePtr.Children.First, modulePtr, nodekind)
		if objPtr != nil {
			return objPtr
		}
	}
	return objPtr
}
