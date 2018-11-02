package serial

import (
	"io"
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
