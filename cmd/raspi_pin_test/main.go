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
	fmt.Println(pin.DigitalRead())
}
