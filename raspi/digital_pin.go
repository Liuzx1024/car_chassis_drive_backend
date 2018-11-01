package raspi

import (
	"errors"
	"sync"
)

const (
	// IN pin mode in
	IN = 0 //"in"
	// OUT pin mode out
	OUT = 1
	// HIGH pin value HIGH
	HIGH = 1
	// LOW pin value LOW
	LOW = 0

	_GPIOClassPath = "/sys/class/gpio"
)

// Errors that are used by the package

var ErrPinNotExported = errors.New("Given pin has not been exported.")
var ErrInvalidPinValue = errors.New("Given value is invalid.")
var ErrInvalidPinMode = errors.New("Given mode is invalid.")

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
		return emptyValue, ErrPinNotExported
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
		return emptyMode, ErrPinNotExported
	}
	return getPinMode(_this.realPin)
}

func (_this *DigitalPin) IsUseAble() bool {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	return _this.useable
}
