package main

import (
	"backend/raspi" //The package will invokes panic if you are not using RaspberryPi or don't have RIGHT permisson!!!
	"fmt"
)

func main() {
	pin, err := raspi.Raspi.ExportPin(3)
	if err != nil {
		panic(err)
	}
	value, err := pin.DigitalRead()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	err = pin.SetPinMode(raspi.OUT)
	if err != nil {
		panic(err)
	}
	err = pin.DigitalWrite(raspi.LOW)
	if err != nil {
		panic(err)
	}
	value, err = pin.DigitalRead()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	fmt.Scanln()
	mode, err := pin.GetPinMode()
	fmt.Println(mode)
	err = raspi.Raspi.UnexportPin(3)
	if err != nil {
		panic(err)
	}
}
