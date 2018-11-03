package serial

import (
	"errors"
	"sync"
	"time"
)

type Serial struct {
	port   *port
	config *config
	mutex  *sync.Mutex
}

type StopBits byte

const (
	Stop1 StopBits = 1
	Stop2 StopBits = 2
)

type Parity byte

const (
	ParityNone Parity = 'N'
	ParityOdd  Parity = 'O'
	ParityEven Parity = 'E'
)

const DefaultSize = 8

func NewSerial(device string) (*Serial, error) {
	conf, err := newDefaultConfig(device)
	if err != nil {
		return nil, err
	}
	obj := &Serial{
		config: conf,
		port:   nil,
		mutex:  new(sync.Mutex),
	}
	return obj, nil
}

func (_this *Serial) SetStopBits(stopBits StopBits) error {
	return _this.changeConfig(func() error {
		return _this.config.setStopBit(stopBits)
	})
}

func (_this *Serial) SetParity(parity Parity) error {
	return _this.changeConfig(func() error {
		return _this.config.setParity(parity)
	})
}

func (_this *Serial) SetSize(size byte) error {
	return _this.changeConfig(func() error {
		return _this.config.setSize(size)
	})
}

func (_this *Serial) SetBaudrate(baud int) error {
	return _this.changeConfig(func() error {
		return _this.config.setBaud(baud)
	})
}

func (_this *Serial) SetReadTimeout(timeout time.Duration) error {
	return _this.changeConfig(func() error {
		return _this.config.setReadTimeout(timeout)
	})
}

func (_this *Serial) Open() error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if err := _this.open(); err != nil {
		return err
	}
	return nil
}

func (_this *Serial) Close() error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if err := _this.close(); err != nil {
		return err
	}
	return nil
}

func (_this *Serial) IsOpened() bool {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	return _this.port != nil
}

var ErrSerialPortNotOpen = errors.New("Serial port is not open")

func (_this *Serial) Write(b []byte) (n int, err error) {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if _this.port != nil {
		return _this.port.write(b)
	}
	return 0, ErrSerialPortNotOpen
}

func (_this *Serial) Read(b []byte) (n int, err error) {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if _this.port != nil {
		return _this.port.read(b)
	}
	return 0, ErrSerialPortNotOpen
}

func (_this *Serial) Flush() error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if _this.port != nil {
		return _this.port.flush()
	}
	return ErrSerialPortNotOpen
}
