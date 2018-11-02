// +build linux

package serial

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func openPort(config *Config) (*Port, error) {
	config.mutex.Lock()
	defer config.mutex.Unlock()
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

	return &Port{file: file}, nil
}
