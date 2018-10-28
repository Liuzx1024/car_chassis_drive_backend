package raspi

import (
	"errors"
)

const (
	IN       = "in"
	OUT      = "out"
	HIGH     = 1
	LOW      = 0
	GPIOPATH = "/sys/class/gpio"
)

var errNotExported = errors.New("pin has not been exported")

type DigitalPin struct {
	pin   string
	label string
}
