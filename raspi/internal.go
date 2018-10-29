package raspi

import (
	"io/ioutil"
	"os"
	"strconv"
)

func generateGpioDirectoryPath(pin uint8) string {
	return GPIOPATH + "gpio" + pinTostring(pin)
}

func pinTostring(pin uint8) string {
	return strconv.Itoa(int(pin))
}

func modeToString(mode uint8) (string, error) {
	switch mode {
	case 0:
		return "in", nil
	case 1:
		return "out", nil
	default:
		return "", ErrInvalidMode
	}
}

// Only use these functions in a critical scope!!!
// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
func exportPin(pin uint8) error {
	exportFile, err := os.OpenFile(GPIOPATH+"/export", os.O_WRONLY, os.ModeType)
	defer func() {
		if exportFile != nil {
			exportFile.Close()
		}
	}()
	if err != nil {
		return err
	}
	_, err = exportFile.WriteString(pinTostring(pin))
	if err != nil {
		return err
	}
	return nil
}

func unexportPin(pin uint8) error {
	unexportFile, err := os.OpenFile(GPIOPATH+"/unexport", os.O_WRONLY, os.ModeType)
	defer func() {
		if unexportFile != nil {
			unexportFile.Close()
		}
	}()
	if err != nil {
		return err
	}
	_, err = unexportFile.WriteString(pinTostring(pin))
	if err != nil {
		return err
	}
	return nil
}

func digitalWrite(pin, value uint8) error {
	if value != HIGH && value != LOW {
		return ErrInvalidValue
	}
	err := ioutil.WriteFile(GPIOPATH+"gpio"+pinTostring(pin)+"/value", []byte(strconv.Itoa(int(value))), os.ModeType)
	if err != nil {
		return err
	}
	return nil
}

func digitalRead(pin uint8) (uint8, error) {
	data, err := ioutil.ReadFile(GPIOPATH + "gpio" + pinTostring(pin) + "/value")
	if err != nil {
		return LOW, err
	}
	intData, err := strconv.Atoi(string(data))
	if err != nil {
		return LOW, err
	}
	return uint8(intData), nil
}

func pinMode(pin uint8, mode uint8) error {
	modeString, err := modeToString(mode)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(GPIOPATH+"gpio"+pinTostring(pin)+"/direction", []byte(modeString), os.ModeType)
	if err != nil {
		return err
	}
	return nil
}

func isPinExported(pin uint8) bool {
	_, err := os.Stat(generateGpioDirectoryPath(pin))
	return os.IsExist(err)
}

func hasRightPermissionToExport() bool {
	exportFile, err := os.OpenFile(GPIOPATH+"/export", os.O_WRONLY, os.ModeType)
	defer func() {
		if exportFile != nil {
			exportFile.Close()
		}
	}()
	return !os.IsPermission(err)
}

func hasRightPermissionToUnexport() bool {
	exportFile, err := os.OpenFile(GPIOPATH+"/unexport", os.O_WRONLY, os.ModeType)
	defer func() {
		if exportFile != nil {
			exportFile.Close()
		}
	}()
	return !os.IsPermission(err)
}
