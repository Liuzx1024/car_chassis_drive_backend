package raspi

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

const emptyString = ""
const cpuinfoFile = "/proc/cpuinfo"

var errRevisonNotFound = errors.New("Can't find revision")
var errNorAValidPin = errors.New("Not a valid pin")
var revision string

func init() {
	if tmpRevision, err := getBoardRevision(); err != nil {
		panic(err)
	} else {
		revision = tmpRevision
		Raspi = new(_raspi)
		Raspi.lock = new(sync.RWMutex)
		Raspi.exportedPin.lock = new(sync.RWMutex)
		Raspi.exportedPin.pins = make(map[uint8]*DigitalPin)
		return
	}
}

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

var Raspi *_raspi

type _raspi struct {
	lock        *sync.RWMutex
	exportedPin *struct {
		lock *sync.RWMutex
		pins map[uint8]*DigitalPin
	}
}

func (_this *_raspi) isRealPinExported(pin uint8) bool {
	_this.exportedPin.lock.RLock()
	defer _this.exportedPin.lock.RUnlock()

	_, ok := _this.exportedPin.pins[pin]
	return ok
}

func (_this *_raspi) ExportPin(pin uint8) (res *DigitalPin, err error) {
	if realPin, pinErr := translatePin(pin); pinErr != nil {
		return nil, errNorAValidPin
	} else {
		_this.exportedPin.lock.Lock()
		defer _this.exportedPin.lock.Unlock()
		if res, ok := _this.exportedPin.pins[pin]; ok {
			return res, nil
		}
		gpioFile, fileErr := os.Open(GPIOPATH + "export")
		defer gpioFile.Close()
		if fileErr != nil {
			return nil, fileErr
		}
		_, fileErr = gpioFile.Write([]byte(strconv.Itoa(int(realPin))))
		if fileErr != nil {
			return nil, fileErr
		}
		gpioDirectoryPath := GPIOPATH + "gpio" + strconv.Itoa(int(realPin))
		if _, fileErr = os.Stat(gpioDirectoryPath); os.IsNotExist(fileErr) {
			return nil, fileErr
		}
		res = &DigitalPin{
			lock:              new(sync.Mutex),
			realPin:           realPin,
			gpioDirectoryPath: gpioDirectoryPath,
		}
		_this.exportedPin.pins[realPin] = res
		return res, nil
	}
}

func (_this *_raspi) UnExportPin(pin uint8) error {
	if realPin, pinErr := translatePin(pin); pinErr != nil {
		return errNorAValidPin
	} else {
		_this.exportedPin.lock.Lock()
		defer _this.exportedPin.lock.Unlock()

		pinObj := _this.exportedPin.pins[realPin]
		pinObj.lock.Lock()
		defer pinObj.lock.Unlock()
		gpioFile, fileErr := os.Open(GPIOPATH + "unexport")
		defer gpioFile.Close()
		if fileErr != nil {
			return fileErr
		}
		_, fileErr = gpioFile.Write([]byte(strconv.Itoa(int(realPin))))
		if fileErr != nil {
			return fileErr
		}
		if _, fileErr = os.Stat(pinObj.gpioDirectoryPath); os.IsExist(fileErr) {
			return fileErr
		}
		delete(_this.exportedPin.pins, realPin)
		return nil
	}
}
