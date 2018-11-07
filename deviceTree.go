package backend

import (
	"backend/uuid"
	"errors"
	"sync"
)

type node struct {
	mutex      *sync.Mutex
	childNodes []uuid.UUID
	parentNode uuid.UUID
	name       string
}

var nodes map[uuid.UUID]*node
var nodesMutex *sync.RWMutex

func init() {
	nodes = make(map[uuid.UUID]*node)
	nodesMutex = new(sync.RWMutex)
	if newUUID, err := uuid.NewV4(); err != nil {
		panic(err)
	} else {
		rootNode := &node{
			mutex:      new(sync.Mutex),
			name:       "root",
			childNodes: []uuid.UUID{},
		}
		nodes[newUUID] = rootNode
	}
}

func NewNode(name string) (uuid.UUID, error) {
	nodesMutex.Lock()
	defer nodesMutex.Unlock()

	for {
		newUUID, err := uuid.NewV4()
		if err != nil {
			return uuid.UUID{}, err
		}
		if _, ok := nodes[newUUID]; ok != false {
			continue
		}
		rootNode, err := FindNodeByName("root")
		ptr := &node{
			childNodes: []uuid.UUID{},
			mutex:      new(sync.Mutex),
			name:       name,
			parentNode: rootNode,
		}
		nodes[newUUID] = ptr
		return newUUID, nil
	}
}

func DeleteNode(uuid uuid.UUID) error {
	nodesMutex.Lock()
	defer nodesMutex.Unlock()

	if ptr, ok := nodes[uuid]; ok == true {
		ptr.mutex.Lock()
		defer ptr.mutex.Unlock()

		if len(ptr.childNodes) != 0 {
			return errors.New("DeleteNode failed")
		}
		delete(nodes, uuid)
		return nil
	}
	return errors.New("Don't have such node")
}

func FindNodeByName(name string) (uuid.UUID, error) {
	nodesMutex.RLock()
	defer nodesMutex.RUnlock()

	for key, ptr := range nodes {
		if ptr.name == name {
			return key, nil
		}
	}
	return uuid.UUID{}, errors.New("Don't have such node")
}

func BindChildToParent(parent uuid.UUID, child uuid.UUID) error {
	if uuid.Equal(parent, child) {
		return errors.New("Parent and child shouldn't have same uuid")
	}
	nodesMutex.Lock()
	defer nodesMutex.Unlock()

	if parentPtr, ok := nodes[parent]; ok {
		if _, ok := nodes[child]; ok {
			parentPtr.mutex.Lock()
			parentPtr.childNodes = append(parentPtr.childNodes, child)
			parentPtr.mutex.Unlock()
			return nil
		}
	}
	return errors.New("Don't have such node")
}

func BindOperationToNode(node uuid.UUID, operationName string, callback func() error) error {}
