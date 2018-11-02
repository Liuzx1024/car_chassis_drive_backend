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

var errVersionNotFound = errors.New("Can't find revision")
var errProcessDontHaveRightPermission = errors.New("Process dont't have right permission")

func getBoardVersion() (string, error) {
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
	return emptyString, errVersionNotFound
}

type raspi struct {
	version      string
	gpioMapMutex *sync.RWMutex
	gpioMap      map[uint8]*DigitalPin
}

//GetBoardVersion
func (_this raspi) GetBoardVersion() string {
	return _this.version
}

//ExportPin
func (_this *raspi) ExportPin(pin uint8) (*DigitalPin, error) {
	_this.gpioMapMutex.Lock()
	defer _this.gpioMapMutex.Unlock()

	// Get real pin number from a given pin
	realPin, err := translatePin(pin)
	if err != nil {
		return nil, err
	}

	// When a real pin is already exported...
	if tmpPtr, ok := _this.gpioMap[pin]; ok {
		// Try to fix this situation:
		// Another application unexported the pin
		if !isPinExported(tmpPtr.realPin) {
			err = exportPin(tmpPtr.realPin)
			if !isPinExported(tmpPtr.realPin) {
				return nil, err
			}
		}
		return tmpPtr, nil
	}

	if !isPinExported(realPin) {
		err := exportPin(realPin)
		if !isPinExported(realPin) {
			return nil, err
		}
	}

	// Construct a Digitalpin structure
	tmpPtr := &DigitalPin{
		lock:    new(sync.Mutex),
		realPin: realPin,
		useable: true,
	}

	// Insert the pointer into the map
	// And return the pointer
	_this.gpioMap[pin] = tmpPtr
	return tmpPtr, nil
}

//UnexportPin
func (_this *raspi) UnexportPin(pin uint8) error {
	_this.gpioMapMutex.Lock()
	defer _this.gpioMapMutex.Unlock()
	if tmpPtr, ok := _this.gpioMap[pin]; !ok {
		return ErrPinNotExported
	} else {
		if isPinExported(tmpPtr.realPin) {
			err := unexportPin(tmpPtr.realPin)
			if isPinExported(tmpPtr.realPin) {
				return err
			}
		}
		tmpPtr.lock.Lock()
		tmpPtr.realPin = emptyPin
		tmpPtr.useable = false
		tmpPtr.lock.Unlock()
		delete(_this.gpioMap, pin)
		return nil
	}
}

// Raspi The Global Raspi pointer is pointed to package's main object
var Raspi *raspi

func init() {
	if !hasRightPermission() {
		panic(errProcessDontHaveRightPermission)
	}
	Raspi = new(raspi)
	if version, err := getBoardVersion(); err != nil {
		panic(err) //If the program is not runnning on RaspberryPi,then invokes panic()
	} else {
		Raspi.version = version
	}
	Raspi.gpioMapMutex = new(sync.RWMutex)
	Raspi.gpioMap = make(map[uint8]*DigitalPin)
}
