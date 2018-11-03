package serialbus

import "errors"

type SerialBus struct {
	m *Master
	s []*Slave
}

var ErrBadPointer = errors.New("Given pointer is nil")

func NewSerialBus(master *Master, slaves ...*Slave) (*SerialBus, error) {
	obj := &SerialBus{
		s: make([]*Slave, len(slaves)),
	}
	if master == nil {
		return nil, ErrBadPointer
	}
	for _, ptr := range slaves {
		if ptr == nil {
			return nil, ErrBadPointer
		}
		obj.s = append(obj.s, ptr)
	}
	return nil, nil
}
