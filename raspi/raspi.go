package raspi

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
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

//GetBoardVersion This function return the version of the RaspberryPi that the program running on.It's used for select a right io layout.
func (_this raspi) GetBoardVersion() string {
	return _this.version
}

//ExportPin This function exports a pin.The "pin" parameter is the pin number of physical io interface on the board.When the opertion is successful,it return an non-nil pointer pointed to a DigitalPin structure stored in a map inside the global raspi structure and a nil error.
func (_this *raspi) ExportPin(pin uint8) (*DigitalPin, error) {
	_this.gpioMapMutex.Lock()
	defer _this.gpioMapMutex.Unlock()

	// Get real pin number from a given pin
	realPin, err := translatePin(pin)
	if err != nil {
		return nil, err
	}

	// When a real pin has been already exported and stored in the map...
	if tmpPtr, ok := _this.gpioMap[pin]; ok {
		// Try to fix this situation:
		// Another application unexported the pin.
		if !isPinExported(tmpPtr.realPin) {
			err = exportPin(tmpPtr.realPin)
			if !isPinExported(tmpPtr.realPin) {
				return nil, err
			}
		}
		return tmpPtr, nil
	}

	//Only call exportPin when it's necessary.
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

//UnexportPin This function export a given pin.The "pin" parameter is the pin number of physical io interface on the board.When the opertion is successful,it return nil.
func (_this *raspi) UnexportPin(pin uint8) error {
	_this.gpioMapMutex.Lock()
	defer _this.gpioMapMutex.Unlock()
	if tmpPtr, ok := _this.gpioMap[pin]; ok {
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
	return ErrPinNotExported
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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		fmt.Println("Raspi is cleanning up...")
		for key := range Raspi.gpioMap {
			Raspi.UnexportPin(key)
		}
		os.Exit(0)
	}()
}
