package gosmi

import (
	"github.com/sleepinggenius2/gosmi/smi"
)

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
	for element := smi.GetFirstElement(n.smiNode); element != nil; element = smi.GetNextElement(element) {
		object := smi.GetElementNode(element)
		if object == nil {
			// TODO: error
			return
		}
		objects = append(objects, CreateNode(object))
	}
	return
}
