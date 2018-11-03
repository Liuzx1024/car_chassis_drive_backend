package serialbus

import (
	"backend/serial"
	"errors"
	"sync"
)

type SerialBus struct {
	master *Master
	slaves []*Slave
	mutex  *sync.Mutex
}

const _FINISHSignal string = "FINISH"

var ErrBadPointer = errors.New("Given pointer is nil")

func (_this *SerialBus) doWork() error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if err := _this.master.doWithSerial(func(serial *serial.Serial) error {
		return serial.Flush()
	}); err != nil {
		return err
	}
	for _, slave := range _this.slaves {
		if err := _this.master.doWithSerial(func(serial *serial.Serial) error {
			return slave.takeTurn(serial)
		}); err == nil {
			return err
		}
	}
	return nil
}

func NewSerialBus(master *Master, slaves ...*Slave) (*SerialBus, error) {
	sb := &SerialBus{
		slaves: []*Slave{},
		mutex:  new(sync.Mutex),
	}
	if master == nil {
		return nil, ErrBadPointer
	}
	for _, slave := range slaves {
		if slave != nil {
			sb.slaves = append(sb.slaves, slave)
		}
	}
	return sb, nil
}
