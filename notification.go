package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

type Notification struct {
	Node
	Objects []Node
}

func (n Node) AsNotification() Notification {
	return Notification{
		Node: n,
		Objects: n.GetNotificationObjects(),
	}
}

func (n Node) GetNotificationObjects() (objects []Node) {
	for element := C.smiGetFirstElement(n.SmiNode); element != nil; element = C.smiGetNextElement(element) {
		object := C.smiGetElementNode(element)
		if object == nil {
			// TODO: error
			return
		}
		objects = append(objects, CreateNode(object))
	}
	return
}
