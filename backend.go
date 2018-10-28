package backend

import (
	"backend/uuid"
)

type AdaptorHandler uuid.UUID
type PlatformHandler uuid.UUID
type DeviceHandler uuid.UUID

func NewAdaptor(name string, driver Driver) (res AdaptorHandler, err error) {
	return
}

func NewPlatform(name string, handler AdaptorHandler, driver Driver) (res PlatformHandler, err error) {
	return
}

func NewDevice(name string, handler PlatformHandler, driver Driver) (res DeviceHandler, err error) {
	return
}
