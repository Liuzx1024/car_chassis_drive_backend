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

func (_this *Slave) readData(r io.Reader) error {
	_this.recvBufMutex.Lock()
	defer _this.recvBufMutex.Unlock()
	pr, pw := io.Pipe()
	go _this.recvBuf.ReadFrom(pr)
	sc := bufio.NewScanner(r)
	sc.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		pw.Write(data)
		if bytes.Contains(data, []byte(_FINISHSignal)) {
			return len(_FINISHSignal), []byte(_FINISHSignal), nil
		} else if atEOF {
			return 0, []byte{}, errFINISHSignalNotFound
		}
		return 0, nil, nil
	})
	for sc.Scan() {
	}
	pw.Close()
	if sc.Err() != nil {
		return errFINISHSignalNotFound
	}
	return nil
}

func (_this *Slave) sendData(w io.Writer) error {
	_this.sendBufMutex.RLock()
	defer _this.sendBufMutex.RUnlock()
	_, err := _this.sendBuf.WriteTo(w)
	if err != nil {
		return err
	}
	return nil
}

func (_this *Slave) sendFINISHSignal(w io.Writer) error {
	_, err := io.WriteString(w, _FINISHSignal)
	return err
}

func (_this *Slave) takeTurn(rw io.ReadWriter) error {
	_this.setCE()
	defer _this.unsetCE()
	if err := _this.sendData(rw); err != nil {
		return err
	}
	if err := _this.sendFINISHSignal(rw); err != nil {
		return err
	}
	if err := _this.readData(rw); err != nil {
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
