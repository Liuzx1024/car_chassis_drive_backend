package backend

import (
	"backend/uuid"
	"container/list"
	"errors"
	"sync"
)

type _globalDeviceTree struct {
	nodeList      *list.List
	nodeMapByUUID map[uuid.UUID]*list.Element
	nodeMapByName map[string]*list.Element
	lock          *sync.RWMutex
}

var globalDeviceTree *_globalDeviceTree

var errGlobalDeviceTreeHasSameUUID = errors.New("errGlobalDeviceTreeHasSameUUID")
var errGlobalDeviceTreeHasSameName = errors.New("errGlobalDeviceTreeHasSameName")
var errGlobalDeviceTreeDontHaveThisNode = errors.New("errGlobalDeviceTreeDontHaveThisNode")

func (_this *_globalDeviceTree) newNode(name string) (*node, error) {
	_this.lock.Lock()
	aNewUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	_, ok := _this.nodeMapByUUID[aNewUUID]
	if ok {
		return nil, errGlobalDeviceTreeHasSameUUID
	}
	_, ok = _this.nodeMapByName[name]
	if ok {
		return nil, errGlobalDeviceTreeHasSameName
	}
	aNewNode := &node{
		lock:     new(sync.RWMutex),
		im:       aNewUUID,
		subNodes: []uuid.UUID{},
		useable:  true,
	}
	aNewElement := _this.nodeList.InsertAfter(aNewNode, _this.nodeList.Front())
	_this.nodeMapByName[name] = aNewElement
	_this.nodeMapByUUID[aNewUUID] = aNewElement
	defer aNewNode.lock.Unlock()
	defer _this.lock.Unlock()
	return aNewNode, nil
}

func (_this *_globalDeviceTree) deleteNode(target *node) error {
	target.lock.Lock()
	_this.lock.Lock()
	defer _this.lock.Unlock()
	defer target.lock.Unlock()
	element, ok := _this.nodeMapByUUID[target.im]
	if !ok {
		return errGlobalDeviceTreeDontHaveThisNode
	}
	_this.nodeList.Remove(element)
	return nil
}

type node struct {
	lock       *sync.RWMutex
	im         uuid.UUID
	parentNode uuid.UUID
	subNodes   []uuid.UUID
	useable    bool
	driver     Driver
}
