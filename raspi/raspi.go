package raspi

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

const (
	emptyString = ""
	emptyMode   = 0
	emptyValue  = 0
	emptyResult = 0
	emptyPin    = 0
	cpuinfoFile = "/proc/cpuinfo"
)

var errRevisonNotFound = errors.New("Can't find revision.")
var errProcessDontHaveRightPermission = errors.New("Process dont't have right permission.")

func getBoardRevision() (string, error) {
	content, err := ioutil.ReadFile(cpuinfoFile)
	if err != nil {
		return emptyString, err
	}
	for _, v := range strings.Split(string(content), "\n") {
		if strings.Contains(v, "Revision") {
			s := strings.Split(string(v), " ")
			version, _ := strconv.ParseInt("0x"+s[len(s)-1], 0, 64)
			if version <= 3 {
				return "1", nil
			} else if version <= 15 {
				return "2", nil
			} else {
				return "3", nil
			}
		}
	}
	return emptyString, errRevisonNotFound
}

type _raspi struct {
	revision     string
	gpioMapMutex *sync.RWMutex
	gpioMap      map[uint8]*DigitalPin
}

func (_this _raspi) GetBoardRevision() string {
	return _this.revision
}

func (_this *_raspi) ExportPin(pin uint8) (*DigitalPin, error) {
	_this.gpioMapMutex.Lock()
	defer _this.gpioMapMutex.Unlock()
	realPin, err := translatePin(pin)
	if err != nil {
		return nil, err
	}

	if tmpPtr, ok := _this.gpioMap[pin]; ok {
		return tmpPtr, nil
	}

	err = exportPin(realPin)
	if err != nil {
		return nil, err
	}

	tmpPtr := &DigitalPin{
		lock:    new(sync.Mutex),
		realPin: realPin,
		useable: true,
	}

	_this.gpioMap[pin] = tmpPtr
	return tmpPtr, nil
}

func (_this *_raspi) UnexportPin(pin uint8) error {
	_this.gpioMapMutex.Lock()
	defer _this.gpioMapMutex.Unlock()
	if tmpPtr, ok := _this.gpioMap[pin]; !ok {
		return ErrPinNotExported
	} else {
		tmpPtr.lock.Lock()
		defer tmpPtr.lock.Unlock()
		tmpPtr.realPin = emptyPin
		tmpPtr.useable = false
		delete(_this.gpioMap, pin)
	}
	return nil
}

func (_this _raspi) GetDigitalPin(pin uint8) (*DigitalPin, error) {
	if tmpPtr, ok := _this.gpioMap[pin]; !ok {
		return nil, ErrPinNotExported
	} else {
		return tmpPtr, nil
	}
}

var Raspi *_raspi

func init() {
	if !hasRightPermissionToExport() || !hasRightPermissionToUnexport() {
		panic(errProcessDontHaveRightPermission)
	}
	Raspi = new(_raspi)
	if revision, err := getBoardRevision(); err != nil {
		panic(err)
	} else {
		Raspi.revision = revision
	}
	Raspi.gpioMapMutex = new(sync.RWMutex)
	Raspi.gpioMap = make(map[uint8]*DigitalPin)
}
