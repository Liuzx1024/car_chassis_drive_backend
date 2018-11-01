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

//Errors that are used by the package

//ErrPinNotExported
var ErrPinNotExported = errors.New("Given pin has not been exported")

//ErrInvalidPinMode
var ErrInvalidPinValue = errors.New("Given value is invalid")

//ErrInvalidPinValue
var ErrInvalidPinMode = errors.New("Given mode is invalid")

//DigitalPin
type DigitalPin struct {
	realPin uint8
	useable bool
	lock    *sync.Mutex
}

//DigitalWrite
func (_this *DigitalPin) DigitalWrite(value uint8) error {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return ErrPinNotExported
		}
	}
	return digitalWrite(_this.realPin, value)
}

//DigitalRead
func (_this *DigitalPin) DigitalRead() (uint8, error) {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return emptyValue, ErrPinNotExported
		}
	}
	return digitalRead(_this.realPin)
}

//SetPinMode
func (_this *DigitalPin) SetPinMode(mode uint8) error {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return ErrPinNotExported
		}
	}
	return setPinMode(_this.realPin, mode)
}

//GetPinMode
func (_this *DigitalPin) GetPinMode() (uint8, error) {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return emptyMode, ErrPinNotExported
		}
	}
	return getPinMode(_this.realPin)
}

//IsUseAble
func (_this *DigitalPin) IsUseAble() bool {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	if _this.useable {
		if !isPinExported(_this.realPin) {
			exportPin(_this.realPin)
			if !isPinExported(_this.realPin) {
				return false
			}
		}
		return true
	}
	return false
}
