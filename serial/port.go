// +build linux

package serial

import (
	"os"

	"golang.org/x/sys/unix"
)

type port struct {
	file *os.File
}

func (_this *port) read(b []byte) (n int, err error) {
	return _this.file.Read(b)
}

func (_this *port) write(b []byte) (n int, err error) {
	return _this.file.Write(b)
}

func (_this *port) flush() error {
	const TCFLSH = 0x540B
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(_this.file.Fd()),
		uintptr(TCFLSH),
		uintptr(unix.TCIOFLUSH),
	)

	if errno == 0 {
		return nil
	}
	return errno
}

func (_this *port) close() (err error) {
	return _this.file.Close()
}
