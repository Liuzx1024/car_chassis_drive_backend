package raspi

import (
	"errors"
	"sync"
)

const (
	IN             = 0 //"in"
	OUT            = 1 //"out"
	HIGH           = 1
	LOW            = 0
	_GPIOClassPath = "/sys/class/gpio/"
)

var ErrPinNotExported = errors.New("Pin has not been exported")
var ErrPinModeNotSatisfy = errors.New("Mode of pin is not satisfy for the request")
var ErrInvalidPinValue = errors.New("Value is invalid")
var ErrInvalidPinMode = errors.New("mode is invalid")

type DigitalPin struct {
	realPin uint8
	useable bool
	lock    *sync.Mutex
}

func (_this *DigitalPin) DigitalWrite(value uint8) error {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		return ErrPinNotExported
	}
	return digitalWrite(_this.realPin, value)
}

func (_this *DigitalPin) DigitalRead() (uint8, error) {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		return 0, ErrPinNotExported
	}
	return digitalRead(_this.realPin)
}

func (_this *DigitalPin) SetPinMode(mode uint8) error {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		return ErrPinNotExported
	}

	return setPinMode(_this.realPin, mode)
}

func (_this *DigitalPin) GetPinMode() (uint8, error) {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	if !isPinExported(_this.realPin) || !_this.useable {
		return 0, ErrPinNotExported
	}
	return getPinMode(_this.realPin)
}
