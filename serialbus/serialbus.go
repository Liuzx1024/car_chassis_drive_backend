package serialbus

import (
	"backend/serial"
	"errors"
	"sync"
)

type SerialBus struct {
	master                     *Master
	slaves                     []*Slave
	mutex                      *sync.Mutex
	workerStatus, workerSignal int
	workerError                error
}

const _FINISHSignal string = "FINISH"

const (
	WorkerRunning = iota
	WorkerStopped
	WorkerExitWithError
	stopWorker
	noSignal
)

var ErrBadPointer = errors.New("Given pointer is nil")

func (_this *SerialBus) doWork(serial *serial.Serial) error {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	for _, slave := range _this.slaves {
		if slave != nil {
			if err := slave.takeTurn(serial); err != nil {
				return err
			}
		}
	}
	return nil
}

func (_this *SerialBus) worker() {
	defer _this.mutex.Unlock()
	for {
		_this.mutex.Lock()
		if err := _this.master.doWithSerial(_this.doWork); err != nil {
			_this.workerError = err
			_this.workerSignal = noSignal
			_this.workerStatus = WorkerExitWithError
		}
		if _this.workerSignal == stopWorker {
			_this.workerError = nil
			_this.workerSignal = noSignal
			_this.workerStatus = WorkerStopped
		}
		_this.mutex.Unlock()
	}
}

func (_this *SerialBus) Run() {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if _this.workerStatus != WorkerRunning {
		_this.workerError = nil
		_this.workerSignal = noSignal
		_this.workerStatus = WorkerRunning
		go _this.worker()
	}
}

func (_this *SerialBus) Stop() {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	if _this.workerStatus == WorkerRunning {
		_this.workerSignal = stopWorker
	}
}

func (_this *SerialBus) Status() (int, error) {
	_this.mutex.Lock()
	defer _this.mutex.Unlock()
	return _this.workerSignal, _this.workerError
}

func NewSerialBus(master *Master, slaves ...*Slave) (*SerialBus, error) {
	newBus := &SerialBus{
		slaves:       []*Slave{},
		mutex:        new(sync.Mutex),
		workerError:  nil,
		workerSignal: noSignal,
		workerStatus: WorkerStopped,
	}
	if master == nil {
		return nil, ErrBadPointer
	}
	for _, slave := range slaves {
		if slave != nil {
			newBus.slaves = append(newBus.slaves, slave)
		}
	}
	return newBus, nil
}
