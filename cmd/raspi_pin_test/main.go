package main

import (
	"backend/raspi" //The package will invokes panic if you are not using RaspberryPi or don't have RIGHT permisson!!!
	"fmt"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pin, err := raspi.Raspi.ExportPin(3)
	panicErr(err)
	value, err := pin.DigitalRead()
	panicErr(err)
	fmt.Println(value)
	err = pin.SetPinMode(raspi.OUT)
	panicErr(err)
	err = pin.DigitalWrite(raspi.LOW)
	panicErr(err)
	value, err = pin.DigitalRead()
	panicErr(err)
	fmt.Println(value)
	fmt.Scanln()
	mode, err := pin.GetPinMode()
	fmt.Println(mode)
	err = pin.DigitalWrite(raspi.HIGH)
	panicErr(err)
	value, err = pin.DigitalRead()
	panicErr(err)
	fmt.Println(value)
	err = raspi.Raspi.UnexportPin(3)
	panicErr(err)
}
