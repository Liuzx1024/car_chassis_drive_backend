package serialbus

import (
	"backend/serial"
	"io"
)

type Master struct {
	rw io.ReadWriter
}

func NewMaster(port *serial.Port) (*Master, error) {
	if port == nil {
		return nil, ErrBadPointer
	}
	obj := &Master{
		rw: port,
	}
	return obj, nil
}
