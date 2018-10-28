package backend

import (
	"backend/uuid"
	"container/list"
	"errors"
	"sync"
)

const emptyString = ""

type _globalDeviceTree struct {
	nodeList      *list.List
	nodeMapByUUID map[uuid.UUID]*list.Element
	lock          *sync.RWMutex
}

var globalDeviceTree *_globalDeviceTree

var errGlobalDeviceTreeHasSameUUID = errors.New("errGlobalDeviceTreeHasSameUUID")
var errGlobalDeviceTreeDontHaveThisNode = errors.New("errGlobalDeviceTreeDontHaveThisNode")

func (_this *_globalDeviceTree) newNode() (*node, error) {
	_this.lock.Lock()
	if aNewUUID, err := uuid.NewV4(); err != nil {
		return nil, err
	} else if _, ok := _this.nodeMapByUUID[aNewUUID]; ok {
		return nil, errGlobalDeviceTreeHasSameUUID
	} else {
		aNewNode := &node{
			lock:     new(sync.RWMutex),
			im:       aNewUUID,
			subNodes: []uuid.UUID{},
			useable:  true,
		}
		aNewElement := _this.nodeList.InsertAfter(aNewNode, _this.nodeList.Back())
		_this.nodeMapByUUID[aNewUUID] = aNewElement
		defer aNewNode.lock.Unlock()
		defer _this.lock.Unlock()
		return aNewNode, nil
	}
}

func (_this *_globalDeviceTree) deleteNode(target *node) error {
	target.lock.Lock()
	_this.lock.Lock()
	defer _this.lock.Unlock()
	defer target.lock.Unlock()
	if element, ok := _this.nodeMapByUUID[target.im]; !ok {
		return errGlobalDeviceTreeDontHaveThisNode
	} else {
		_this.nodeList.Remove(element)
		delete(_this.nodeMapByUUID, target.im)
		target.useable = false
	}
	return nil
}

func (_this *_globalDeviceTree) lookupNodeByUUID(uuid uuid.UUID) (res *node, err error) {
	_this.lock.RLock()
	defer _this.lock.RUnlock()
	if tmp, ok := _this.nodeMapByUUID[uuid]; ok {
		return nil, errGlobalDeviceTreeDontHaveThisNode
	} else {
		return tmp.Value.(*node), nil
	}
}

type node struct {
	lock       *sync.RWMutex
	im         uuid.UUID
	parentNode uuid.UUID
	subNodes   []uuid.UUID
	useable    bool
	driver     Driver
	name       string
}

func (_this *node) check() bool {
	return _this.useable
}
