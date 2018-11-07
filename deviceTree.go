package backend

import (
	"backend/uuid"
	"sync"
)

type node struct {
}

type _nodes struct {
	nodes map[uuid.UUID]node
	mutex *sync.Mutex
}

func (_this *_nodes) NewNode() uuid.UUID {
	return uuid.UUID{}
}

func init() {

}
