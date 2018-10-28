package backend

import (
	"backend/uuid"
	"container/list"
	"sync"
)

type adaptor node
type controller node
type sensor node
type actuator node

func init() {
	func() {
		globalDeviceTree = new(_globalDeviceTree)
		globalDeviceTree.lock = new(sync.RWMutex)
		globalDeviceTree.lock.Lock()
		globalDeviceTree.nodeList = new(list.List).Init()
		globalDeviceTree.nodeMapByName = make(map[string]*list.Element)
		globalDeviceTree.nodeMapByUUID = make(map[uuid.UUID]*list.Element)
	}()
}
