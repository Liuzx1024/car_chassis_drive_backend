package raspi

import (
	"errors"
	"sync"
)

const (
	//IN       = "in"
	//OUT      = "out"
	IN       = 0
	OUT      = 1
	HIGH     = 1
	LOW      = 0
	GPIOPATH = "/sys/class/gpio/"
)

var ErrNotExported = errors.New("Pin has not been exported")
var ErrPinModeNotSatisfy = errors.New("Mode of pin is not satisfy for the request")
var ErrInvalidValue = errors.New("Value is invalid")
var ErrInvalidMode = errors.New("mode is invalid")

type DigitalPin struct {
	gpioDirectoryPath string
	realPin           uint8
	lock              *sync.Mutex
}
