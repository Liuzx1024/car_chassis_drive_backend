package serialbus

import (
	"backend/raspi"
	"bufio"
	"bytes"
	"errors"
	"io"
	"sync"
)

type Slave struct {
	ce                         *raspi.DigitalPin
	recvBufMutex, sendBufMutex *sync.RWMutex
	recvBuf, sendBuf           *bytes.Buffer
}

func (_this *Slave) setCE() error {
	return _this.ce.DigitalWrite(raspi.LOW)
}

func (_this *Slave) unsetCE() error {
	return _this.ce.DigitalWrite(raspi.HIGH)
}

var errFINISHSignalNotFound = errors.New("FINISH signal not found")

func (_this *Slave) recvData(r io.Reader) error {
	_this.recvBufMutex.Lock()
	defer _this.recvBufMutex.Unlock()
	reader := bufio.NewReader(r)
	for {
		buf, _, err := reader.ReadLine()
		if err != nil {
			return err
		} else {
			if string(buf) == _FINISHSignal {
				break
			}
			_this.recvBuf.Write(buf)
		}
	}
	return nil
}

func (_this *Slave) sendData(w io.Writer) error {
	_this.sendBufMutex.RLock()
	defer _this.sendBufMutex.RUnlock()
	if _, err := _this.sendBuf.WriteTo(w); err != nil {
		return err
	}
	if _, err := io.WriteString(w, _FINISHSignal); err != nil {
		return err
	}
	return nil
}

func (_this *Slave) takeTurn(rw io.ReadWriter) error {
	_this.setCE()
	defer _this.unsetCE()
	sendErr := make(chan error)
	go func() {
		err := _this.sendData(rw)
		sendErr <- err
	}()
	recvErr := make(chan error)
	go func() {
		err := _this.recvData(rw)
		recvErr <- err
	}()
	if err := <-sendErr; err != nil {
		return err
	} else if err = <-recvErr; err != nil {
		return err
	}
	return nil
}

func NewSlave(ce *raspi.DigitalPin) (*Slave, error) {
	if ce == nil {
		return nil, ErrBadPointer
	}
	if err := ce.SetPinMode(raspi.OUTPUT); err != nil {
		return nil, err
	}
	if err := ce.DigitalWrite(raspi.HIGH); err != nil {
		return nil, err
	}
	obj := &Slave{
		ce:           ce,
		recvBuf:      bytes.NewBuffer([]byte{}),
		sendBuf:      bytes.NewBuffer([]byte{}),
		recvBufMutex: new(sync.RWMutex),
		sendBufMutex: new(sync.RWMutex),
	}
	return obj, nil
}
