package serial

import (
	"errors"
	"os"
	"time"

	"golang.org/x/sys/unix"
)

type config struct {
	name        string
	baud        int
	readTimeout time.Duration
	size        byte
	parity      Parity
	stopBits    StopBits
}

func newDefaultConfig(name string) (*config, error) {
	if file, err := os.OpenFile(name, os.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0666); err != nil {
		return nil, err
	} else {
		file.Close()
	}
	return &config{
		name:     name,
		size:     DefaultSize,
		parity:   ParityNone,
		stopBits: Stop1,
		baud:     115200,
	}, nil
}

var bauds = map[int]uint32{
	50:      unix.B50,
	75:      unix.B75,
	110:     unix.B110,
	134:     unix.B134,
	150:     unix.B150,
	200:     unix.B200,
	300:     unix.B300,
	600:     unix.B600,
	1200:    unix.B1200,
	1800:    unix.B1800,
	2400:    unix.B2400,
	4800:    unix.B4800,
	9600:    unix.B9600,
	19200:   unix.B19200,
	38400:   unix.B38400,
	57600:   unix.B57600,
	115200:  unix.B115200,
	230400:  unix.B230400,
	460800:  unix.B460800,
	500000:  unix.B500000,
	576000:  unix.B576000,
	921600:  unix.B921600,
	1000000: unix.B1000000,
	1152000: unix.B1152000,
	1500000: unix.B1500000,
	2000000: unix.B2000000,
	2500000: unix.B2500000,
	3000000: unix.B3000000,
	3500000: unix.B3500000,
	4000000: unix.B4000000,
}

func (_this *config) setReadTimeout(time time.Duration) error {
	_this.readTimeout = time
	return nil
}

var ErrBadBaud = errors.New("unsuportted baud setting")

func (_this *config) setBaud(baud int) error {
	if _, ok := bauds[baud]; ok != true {
		return ErrBadBaud
	}
	_this.baud = baud
	return nil
}

var ErrBadSize error = errors.New("unsupported serial data size")

func (_this *config) setSize(size byte) error {
	if size < 5 || size > 8 {
		return ErrBadBaud
	}
	_this.size = size
	return nil
}

var ErrBadStopBits error = errors.New("unsupported stop bit setting")

func (_this *config) setStopBit(stopBits StopBits) error {
	if stopBits != Stop1 && stopBits != Stop2 {
		return ErrBadStopBits
	}
	_this.stopBits = stopBits
	return nil
}

var ErrBadParity error = errors.New("unsupported parity setting")

func (_this *config) setParity(parity Parity) error {
	if parity != ParityNone && parity != ParityEven && parity != ParityOdd {
		return ErrBadParity
	}
	_this.parity = parity
	return nil
}
