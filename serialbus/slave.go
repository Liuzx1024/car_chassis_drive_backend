package serialbus

import (
	"backend/raspi"
	"bufio"
	"bytes"
	"encoding/json"
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
			_this.recvBuf.Write([]byte("\n"))
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
	wg := new(sync.WaitGroup)
	wg.Add(2)
	sendErr := make(chan error)
	go func() {
		err := _this.sendData(rw)
		sendErr <- err
		wg.Done()
	}()
	recvErr := make(chan error)
	go func() {
		err := _this.recvData(rw)
		recvErr <- err
		wg.Done()
	}()
	wg.Wait()
	if err := <-sendErr; err != nil {
		return err
	} else if err = <-recvErr; err != nil {
		return err
	}
	return nil
}

func (_this *Slave) StoreJSONMessageToSendBuf(message json.RawMessage) error {
	_this.sendBufMutex.Lock()
	defer _this.sendBufMutex.Unlock()
	_this.sendBuf.Write(message)
	_this.sendBuf.Write([]byte("\n"))
	return nil
}

func (_this *Slave) LoadJSONMessageFromRecvBuf() (json.RawMessage, error) {
	_this.recvBufMutex.RLock()
	defer _this.recvBufMutex.RUnlock()
	raw, _, err := bufio.NewReader(_this.recvBuf).ReadLine()
	return json.RawMessage(raw), err
}

func NewSlave(pin uint8) (*Slave, error) {
	ce, err := raspi.ExportPin(pin)
	if err != nil {
		return nil, err
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
