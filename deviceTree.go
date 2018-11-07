package backend

import (
	"backend/uuid"
	"errors"
	"sync"
)

type node struct {
	mutex      *sync.Mutex
	childNodes []uuid.UUID
}

var nodes map[uuid.UUID]*node
var nodesMutex *sync.RWMutex

func init() {
	nodes = make(map[uuid.UUID]*node)
	nodesMutex = new(sync.RWMutex)
}

func NewNode() (uuid.UUID, error) {
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
		ptr := &node{
			childNodes: []uuid.UUID{},
			mutex:      new(sync.Mutex),
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
