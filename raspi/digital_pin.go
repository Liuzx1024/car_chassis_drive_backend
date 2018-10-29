package raspi

import (
	"errors"
	"sync"
)

const (
	IN       = "in"
	OUT      = "out"
	HIGH     = 1
	LOW      = 0
	GPIOPATH = "/sys/class/gpio/"
)

var errNotExported = errors.New("pin has not been exported")

type DigitalPin struct {
	gpioDirectoryPath string
	realPin           uint8
	lock              *sync.Mutex
}

func (_this *DigitalPin) check() bool {
	return Raspi.isRealPinExported(_this.realPin)
}

func (_this *DigitalPin) DigitalWrite(pin, value uint8) error {
	if !_this.check() {
		return errNotExported
	} else {
		return nil
	}
}

func (_this *DigitalPin) DigitalRead(pin uint8) error {
	if !_this.check() {
		return errNotExported
	} else {
		return nil
	}
}

func (_this *DigitalPin) PinMode(pin uint8, mode uint8) error {
	if !_this.check() {
		return errNotExported
	} else {
		return nil
	}
}
