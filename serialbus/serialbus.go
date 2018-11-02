package serialbus

import (
	"errors"
	"sync/atomic"
	"time"
)

type SerialBus struct {
	m            *Master
	s            []*Slave
	workerStatus int64
}

const _FINISHSignal string = "FINISH\n"

var ErrBadPointer = errors.New("Given pointer is nil")
var ErrMasterHasBeenOpened = errors.New("Given worker has been opened")

const (
	workerIsNotRunning = iota
	workerIsRunnning
)

const defaultDelayTime = time.Millisecond

func (_this *SerialBus) worker() {
	defer atomic.StoreInt64(&_this.workerStatus, workerIsNotRunning)
	_this.m.open()
	defer _this.m.close()
	if err := _this.m.open(); err != nil {
		return
	}
	defer _this.m.close()
	timer := time.NewTimer(defaultDelayTime)
	for {
		timer.Reset(defaultDelayTime)
		_this.m.mutex.RLock()
		for _, ptr := range _this.s {
			if ptr == nil {
				panic(ErrBadPointer)
			}
			//here every slave take its turn
			if err := ptr.takeTurn(_this.m.port); err != nil {
				_this.m.mutex.RUnlock()
				return
			}
			_this.m.flush()
		}
		_this.m.mutex.RUnlock()
		<-timer.C
		if atomic.LoadInt64(&_this.workerStatus) == workerIsNotRunning {
			return
		}
	}
}

func (_this *SerialBus) StartWorker() {
	if atomic.LoadInt64(&_this.workerStatus) == workerIsNotRunning {
		atomic.StoreInt64(&_this.workerStatus, workerIsRunnning)
		go _this.worker()
	}
}

func (_this *SerialBus) StopWorker() {
	if atomic.LoadInt64(&_this.workerStatus) == workerIsRunnning {
		atomic.StoreInt64(&_this.workerStatus, workerIsNotRunning)
	}
}

func (_this *SerialBus) GetWorkerStatus() int64 {
	return atomic.LoadInt64(&_this.workerStatus)
}

func NewSerialBus(master *Master, slaves ...*Slave) (*SerialBus, error) {
	obj := &SerialBus{
		s: make([]*Slave, len(slaves)),
	}
	if master == nil {
		return nil, ErrBadPointer
	}
	if master.port != nil {
		return nil, ErrMasterHasBeenOpened
	}
	for _, ptr := range slaves {
		if ptr == nil {
			return nil, ErrBadPointer
		}
		obj.s = append(obj.s, ptr)
	}
	atomic.StoreInt64(&obj.workerStatus, workerIsNotRunning)
	return obj, nil
}
