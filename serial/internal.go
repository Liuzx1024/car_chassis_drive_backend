// +build linux

package serial

import (
	"os"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func openPort(config *config) (*port, error) {
	file, err := os.OpenFile(config.name, unix.O_RDWR|unix.O_NOCTTY|unix.O_NONBLOCK, 0666)
	if err != nil {
		if file != nil {
			file.Close()
		}
		return nil, err
	}

	rate, ok := bauds[config.baud]

	if !ok {
		return nil, ErrBadBaud
	}

	cflagToUse := unix.CREAD | unix.CLOCAL | rate

	switch config.size {
	case 5:
		cflagToUse |= unix.CS5
	case 6:
		cflagToUse |= unix.CS6
	case 7:
		cflagToUse |= unix.CS7
	case 8:
		cflagToUse |= unix.CS8
	default:
		return nil, ErrBadSize
	}

	switch config.stopBits {
	case Stop1:
		//do nothing
	case Stop2:
		cflagToUse |= unix.CSTOPB
	default:
		return nil, ErrBadStopBits
	}

	switch config.parity {
	case ParityNone:
		//do nothing
	case ParityOdd:
		cflagToUse |= unix.PARENB
		cflagToUse |= unix.PARODD
	case ParityEven:
		cflagToUse |= unix.PARENB
	default:
		return nil, ErrBadParity
	}

	fd := file.Fd()
	vmin, vtime := posixTimeoutValues(config.readTimeout)
	t := unix.Termios{
		Iflag:  unix.IGNPAR,
		Cflag:  cflagToUse,
		Ispeed: rate,
		Ospeed: rate,
	}

	t.Cc[unix.VMIN] = vmin
	t.Cc[unix.VTIME] = vtime

	if _, _, errno := unix.Syscall6(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(unix.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	); errno != 0 {
		return nil, errno
	}

	if err = unix.SetNonblock(int(fd), false); err != nil {
		return nil, err
	}

	return &port{file: file}, nil
}

func posixTimeoutValues(readTimeout time.Duration) (vmin uint8, vtime uint8) {
	const MAXUINT8 = 1<<8 - 1
	var minBytesToRead uint8 = 1
	var readTimeoutInDeci int64
	if readTimeout > 0 {
		minBytesToRead = 0
		readTimeoutInDeci = (readTimeout.Nanoseconds() / 1e6 / 100)
		if readTimeoutInDeci < 1 {
			readTimeoutInDeci = 1
		} else if readTimeoutInDeci > MAXUINT8 {
			readTimeoutInDeci = MAXUINT8
		}
	}
	return minBytesToRead, uint8(readTimeoutInDeci)
}

func (_this *Serial) changeConfig(callback func() error) error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	shouldRestart := false
	if _this.port != nil {
		shouldRestart = true
		if err := _this.close(); err != nil {
			return err
		}
	}
	if err := callback(); err != nil {
		return err
	}
	if shouldRestart {
		if err := _this.open(); err != nil {
			return err
		}
	}
	return nil
}

func (_this *Serial) open() error {
	if _this.port == nil {
		if port, err := openPort(_this.config); err != nil {
			return err
		} else {
			_this.port = port
		}
	}
	return nil
}

func (_this *Serial) close() error {
	if _this.port != nil {
		err := _this.port.close()
		_this.port = nil
		if err != nil {
			return err
		}
	}
	return nil
}
