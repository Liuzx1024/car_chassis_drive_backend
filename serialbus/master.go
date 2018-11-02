package serialbus

import (
	"backend/serial"
	"sync"
)

type Master struct {
	config serial.Config
	port   *serial.Port
	mutex  *sync.RWMutex
}

func (_this *Master) close() error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if _this.port != nil {
		err := _this.port.Close()
		if err != nil {
			return err
		}
		_this.port = nil
	}
	return nil
}

func (_this *Master) flush() error {
	_this.mutex.RLock()
	defer _this.mutex.RUnlock()
	return _this.port.Flush()
}

func (_this *Master) open() error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if _this.port == nil {
		port, err := serial.OpenPort(_this.config)
		if err != nil {
			return err
		}
		_this.port = port
	}
	return nil
}

func NewMaster(config serial.Config) (*Master, error) {
	port, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}
	obj := &Master{
		port:   port,
		config: config,
		mutex:  new(sync.RWMutex),
	}
	if err := obj.close(); err != nil {
		return nil, err
	}
	return obj, nil
}
