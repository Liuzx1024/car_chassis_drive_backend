package raspi

import (
	"errors"
	"io/ioutil"
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

var Raspi _raspi

type _raspi struct {
	lock        sync.RWMutex
	digitalPins map[int]*DigitalPin
}

func (_this *_raspi) ExportDigitalPin(pin string) *DigitalPin {
	return nil
}
