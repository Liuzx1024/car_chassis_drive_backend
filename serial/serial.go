package serial

import (
	"errors"
	"io"
	"math"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type ParityMode int

const (
	PARITY_NONE ParityMode = 0
	PARITY_ODD  ParityMode = 1
	PARITY_EVEN ParityMode = 2
)

var (
	StandardBaudRates = map[uint]bool{
		50:     true,
		75:     true,
		110:    true,
		134:    true,
		150:    true,
		200:    true,
		300:    true,
		600:    true,
		1200:   true,
		1800:   true,
		2400:   true,
		4800:   true,
		7200:   true,
		9600:   true,
		14400:  true,
		19200:  true,
		28800:  true,
		38400:  true,
		57600:  true,
		76800:  true,
		115200: true,
		230400: true,
	}
)

func IsStandardBaudRate(baudRate uint) bool { return StandardBaudRates[baudRate] }

type OpenOptions struct {
	PortName                string
	BaudRate                uint
	DataBits                uint
	StopBits                uint
	ParityMode              ParityMode
	RTSCTSFlowControl       bool
	InterCharacterTimeout   uint
	MinimumReadSize         uint
	Rs485Enable             bool
	Rs485RtsHighDuringSend  bool
	Rs485RtsHighAfterSend   bool
	Rs485RxDuringTx         bool
	Rs485DelayRtsBeforeSend int
	Rs485DelayRtsAfterSend  int
}

func Open(options OpenOptions) (io.ReadWriteCloser, error) {
	return openInternal(options)
}

func round(f float64) float64 {
	return math.Floor(f + 0.5)
}

// Grab the constants with the following little program, to avoid using cgo:
// #include <stdio.h>
// #include <stdlib.h>
// #include <linux/termios.h>
//
// int main(int argc, const char **argv) {
//   printf("kTCSETS2 = 0x%08X\n", TCSETS2);
//   printf("kBOTHER  = 0x%08X\n", BOTHER);
//   printf("kNCCS    = %d\n",     NCCS);
//   return 0;
// }
const (
	kTCSETS2 = 0x402C542B
	kBOTHER  = 0x00001000
	kNCCS    = 19
)

//Fucking asshole-unistd
type cc_t byte
type speed_t uint32
type tcflag_t uint32
type termios2 struct {
	c_iflag  tcflag_t    // input mode flags
	c_oflag  tcflag_t    // output mode flags
	c_cflag  tcflag_t    // control mode flags
	c_lflag  tcflag_t    // local mode flags
	c_line   cc_t        // line discipline
	c_cc     [kNCCS]cc_t // control characters
	c_ispeed speed_t     // input speed
	c_ospeed speed_t     // output speed
}

const (
	sER_RS485_ENABLED        = (1 << 0)
	sER_RS485_RTS_ON_SEND    = (1 << 1)
	sER_RS485_RTS_AFTER_SEND = (1 << 2)
	sER_RS485_RX_DURING_TX   = (1 << 4)
	tIOCSRS485               = 0x542F
)

type serial_rs485 struct {
	flags                 uint32
	delay_rts_before_send uint32
	delay_rts_after_send  uint32
	padding               [5]uint32
}

func makeTermios2(options OpenOptions) (*termios2, error) {
	vtime := uint(round(float64(options.InterCharacterTimeout)/100.0) * 100)
	vmin := options.MinimumReadSize

	if vmin == 0 && vtime < 100 {
		return nil, errors.New("invalid values for InterCharacterTimeout and MinimumReadSize")
	}

	if vtime > 25500 {
		return nil, errors.New("invalid value for InterCharacterTimeout")
	}

	ccOpts := [kNCCS]cc_t{}
	ccOpts[syscall.VTIME] = cc_t(vtime / 100)
	ccOpts[syscall.VMIN] = cc_t(vmin)

	t2 := &termios2{
		c_cflag:  syscall.CLOCAL | syscall.CREAD | kBOTHER,
		c_ispeed: speed_t(options.BaudRate),
		c_ospeed: speed_t(options.BaudRate),
		c_cc:     ccOpts,
	}

	switch options.StopBits {
	case 1:
	case 2:
		t2.c_cflag |= syscall.CSTOPB

	default:
		return nil, errors.New("invalid setting for StopBits")
	}

	switch options.ParityMode {
	case PARITY_NONE:
	case PARITY_ODD:
		t2.c_cflag |= syscall.PARENB
		t2.c_cflag |= syscall.PARODD

	case PARITY_EVEN:
		t2.c_cflag |= syscall.PARENB

	default:
		return nil, errors.New("invalid setting for ParityMode")
	}

	switch options.DataBits {
	case 5:
		t2.c_cflag |= syscall.CS5
	case 6:
		t2.c_cflag |= syscall.CS6
	case 7:
		t2.c_cflag |= syscall.CS7
	case 8:
		t2.c_cflag |= syscall.CS8
	default:
		return nil, errors.New("invalid setting for DataBits")
	}

	if options.RTSCTSFlowControl {
		t2.c_cflag |= unix.CRTSCTS
	}

	return t2, nil
}

func openInternal(options OpenOptions) (io.ReadWriteCloser, error) {

	file, openErr :=
		os.OpenFile(
			options.PortName,
			syscall.O_RDWR|syscall.O_NOCTTY|syscall.O_NONBLOCK,
			0600)
	if openErr != nil {
		return nil, openErr
	}
	nonblockErr := syscall.SetNonblock(int(file.Fd()), false)
	if nonblockErr != nil {
		return nil, nonblockErr
	}

	t2, optErr := makeTermios2(options)
	if optErr != nil {
		return nil, optErr
	}

	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(kTCSETS2),
		uintptr(unsafe.Pointer(t2)))

	if errno != 0 {
		return nil, os.NewSyscallError("SYS_IOCTL", errno)
	}

	if r != 0 {
		return nil, errors.New("unknown error from SYS_IOCTL")
	}

	if options.Rs485Enable {
		rs485 := serial_rs485{
			sER_RS485_ENABLED,
			uint32(options.Rs485DelayRtsBeforeSend),
			uint32(options.Rs485DelayRtsAfterSend),
			[5]uint32{0, 0, 0, 0, 0},
		}

		if options.Rs485RtsHighDuringSend {
			rs485.flags |= sER_RS485_RTS_ON_SEND
		}

		if options.Rs485RtsHighAfterSend {
			rs485.flags |= sER_RS485_RTS_AFTER_SEND
		}

		r, _, errno := syscall.Syscall(
			syscall.SYS_IOCTL,
			uintptr(file.Fd()),
			uintptr(tIOCSRS485),
			uintptr(unsafe.Pointer(&rs485)))

		if errno != 0 {
			return nil, os.NewSyscallError("SYS_IOCTL (RS485)", errno)
		}

		if r != 0 {
			return nil, errors.New("Unknown error from SYS_IOCTL (RS485)")
		}
	}

	return file, nil
}
