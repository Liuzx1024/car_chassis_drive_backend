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
	pin       string
	realPin   int
	direction string
	value     int
}

func newDigitalPin(pin string) (res *DigitalPin, err error) {
	if tmpRes, tmpErr := translatePin(pin); tmpErr != nil {
		return nil, tmpErr
	} else {
		return &DigitalPin{
			pin:     pin,
			realPin: tmpRes,
		}, nil
	}
}
