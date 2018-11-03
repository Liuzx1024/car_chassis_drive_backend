package serialbus

import (
	"backend/serial"
	"sync"
)

type Master struct {
	serial *serial.Serial
	mutex  *sync.Mutex
}

func (_this *Master) doWithSerial(cb func(serial *serial.Serial) error) error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	return cb(_this.serial)
}

func NewMaster(device string, buadRate int) (*Master, error) {
	serial, err := serial.NewSerial(device)
	if err != nil {
		return nil, err
	}
	if err := serial.SetBaudrate(buadRate); err != nil {
		return nil, err
	}
	obj := &Master{
		mutex:  new(sync.Mutex),
		serial: serial,
	}
	return obj, nil
}
