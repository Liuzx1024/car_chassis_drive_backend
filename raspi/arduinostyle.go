package raspi

//ArduinoInterface Just a Interface,pleasw use Raspi if your project is not single-routine.
var ArduinoInterface arduino

type arduino struct{}

//DigitalWrite Just like you are using Arduino.
func (arduino) DigitalWrite(pin, value uint8) {
	ptr, err := Raspi.ExportPin(pin)
	if err != nil {
		return
	}
	ptr.DigitalWrite(value)
}

//DigitalRead Just like you are using Arduino.
func (arduino) DigitalRead(pin uint8) uint8 {
	ptr, err := Raspi.ExportPin(pin)
	if err != nil {
		return emptyValue
	}
	value, err := ptr.DigitalRead()
	if err != nil {
		return emptyValue
	}
	return value
}

//PinMode Just like you are using Arduino,but it return the pin mode after settting it.
func (arduino) PinMode(pin, mode uint8) uint8 {
	ptr, err := Raspi.ExportPin(pin)
	if err != nil {
		return emptyMode
	}
	err = ptr.SetPinMode(mode)
	if err != nil {
		return emptyMode
	}
	mode, err = ptr.GetPinMode()
	if err != nil {
		return emptyMode
	}
	return mode
}
