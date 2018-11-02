package main

import (
	"backend/raspi"
	"time"
)

func dealWithErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pin, err := raspi.Raspi.ExportPin(3)
	dealWithErr(err)
	err = pin.SetPinMode(raspi.OUT)
	dealWithErr(err)
	for {
		err = pin.DigitalWrite(raspi.HIGH)
		dealWithErr(err)
		time.Sleep(time.Second)
		err = pin.DigitalWrite(raspi.LOW)
		dealWithErr(err)
		time.Sleep(time.Second)
	}

}
