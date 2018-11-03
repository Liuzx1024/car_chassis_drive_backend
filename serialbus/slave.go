package serialbus

import (
	"backend/raspi"
)

type Slave struct {
	ce *raspi.DigitalPin
}

func NewSlave(ce *raspi.DigitalPin) (*Slave, error) {
	if ce == nil {
		return nil, ErrBadPointer
	}
	if err := ce.SetPinMode(raspi.OUTPUT); err != nil {
		return nil, err
	}
	if err := ce.DigitalWrite(raspi.HIGH); err != nil {
		return nil, err
	}
	obj := &Slave{
		ce: ce,
	}
	return obj, nil
}
