package internal

import (
	"github.com/sleepinggenius2/gosmi/types"
)

type Node struct {
	SubId       types.SmiSubId
	Flags       Flags
	OidLen      int
	Oid         types.Oid
	Parent      *Node
	Prev        *Node
	Next        *Node
	Children    NodeChildMap
	FirstObject *Object
	LastObject  *Object
}

func (x *Node) AddObject(obj *Object) {
	obj.Node = x
	obj.PrevSameNode = x.LastObject
	if x.LastObject == nil {
		x.FirstObject = obj
	} else {
		x.LastObject.NextSameNode = obj
	}
	x.LastObject = obj
}

func (x *Node) IsRoot() bool {
	return x != nil && x.Flags.Has(FlagRoot)
}

type NodeChildMap struct {
	First *Node

	last *Node
	m    map[types.SmiSubId]*Node
}

func (x *NodeChildMap) Add(n *Node) {
	existing := x.Get(n.SubId)
	if existing != nil {
		for obj := n.FirstObject; obj != nil; obj = obj.NextSameNode {
			existing.AddObject(obj)
		}
		return
	}
	if n.Parent != nil && n.Parent.Oid != nil {
		n.Oid = types.NewOid(n.Parent.Oid, n.SubId)
		n.OidLen = n.Parent.OidLen + 1
	}
	if x.last == nil {
		x.First = n
		x.last = n
	} else {
		c := x.First
		for c != nil && c.SubId < n.SubId {
			c = c.Next
		}
		if c == nil {
			n.Prev = x.last
			x.last.Next = n
			x.last = n
		} else {
			n.Prev = c.Prev
			n.Next = c
			if c.Prev == nil {
				x.First = n
			} else {
				c.Prev.Next = n
			}
			c.Prev = n
		}
	}

	if x.m == nil {
		x.m = make(map[types.SmiSubId]*Node)
	}
	x.m[n.SubId] = n
}

func (x *NodeChildMap) Get(id types.SmiSubId) *Node {
	if x.m == nil {
		return nil
	}
	return x.m[id]
}

func FindNodeByOid(oidlen int, oid types.Oid) *Node {
	nodePtr := smiHandle.RootNode
	for i := 0; i < oidlen && nodePtr != nil; i++ {
		nodePtr = nodePtr.Children.Get(oid[i])
	}
	return nodePtr
}

/*func createNodes(oid types.Oid) *Node {
	var nodePtr *Node
	parentNodePtr := smiHandle.RootNode
	for i := 0; i < len(oid); i++ {
		nodePtr = parentNodePtr.Children.Get(oid[i])
		if nodePtr != nil {
			continue
		}
		nodePtr = &Node{
			SubId:  oid[i],
			Oid:    oid[:i+1],
			Parent: parentNodePtr,
		}
		parentNodePtr.Children.Add(nodePtr)
		parentNodePtr = nodePtr
	}
	return parentNodePtr
}*/
