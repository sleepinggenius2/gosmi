package gosmi

/*
#cgo LDFLAGS: -lsmi
#include <stdlib.h>
#include <smi.h>
*/
import "C"

type Notification struct {
	SmiNode
	Objects []SmiNode
}

func (n SmiNode) AsNotification() Notification {
	return Notification{
		SmiNode: n,
		Objects: n.GetNotificationObjects(),
	}
}

func (n SmiNode) GetNotificationObjects() (objects []SmiNode) {
	for element := C.smiGetFirstElement(n.smiNode); element != nil; element = C.smiGetNextElement(element) {
		object := C.smiGetElementNode(element)
		if object == nil {
			// TODO: error
			return
		}
		objects = append(objects, CreateNode(object))
	}
	return
}
