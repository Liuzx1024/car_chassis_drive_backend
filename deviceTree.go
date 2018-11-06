package backend

import (
	"sync"
)

var nodes *sync.Map

type node struct{}
