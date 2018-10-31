package serialbus

import (
	"backend/raspi"
	"bytes"
	"io"
	"sync"
)

type Slave struct {
	ce                         *raspi.DigitalPin
	recvBufMutex, sendBufMutex *sync.Mutex
	recvBuf, sendBuf           *bytes.Buffer
}

func (_this *Slave) setCE() error {
	return _this.ce.DigitalWrite(raspi.LOW)
}

func (_this *Slave) unsetCE() error {
	return _this.ce.DigitalWrite(raspi.HIGH)
}

func (_this *Slave) readData(Reader io.Reader) {}

func (_this *Slave) sendData(w io.Writer) {}

func (_this *Slave) sendFINISHSignal(w io.Writer) error {
	_, err := io.WriteString(w, _FINISHSignal)
	return err
}

func (_this *Slave) takeTurn(rw io.ReadWriter) {
	_this.setCE()
	_this.sendData(rw)
	_this.sendFINISHSignal(rw)
	_this.readData(rw)
	_this.unsetCE()
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
		recvBufMutex: new(sync.Mutex),
		sendBufMutex: new(sync.Mutex),
	}
	return obj, nil
}
