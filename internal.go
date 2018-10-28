package backend

import (
	"backend/uuid"
	"fmt"
	"sync"
)

type internalAdaptor struct {
	myUUID    uuid.UUID
	platforms map[uuid.UUID]internalPlatform
}

type internalPlatform struct {
	myUUID   uuid.UUID
	devices  map[uuid.UUID]internalDevice
	myParent *internalAdaptor
}

type internalDevice struct {
	myUUID   uuid.UUID
	myParent *internalPlatform
}

type adaptorName string
type platformName string
type deviceName string

type _deviceTree struct {
	lock         sync.RWMutex
	adaptorTree  map[adaptorName]*internalAdaptor
	platformTree map[adaptorName]map[platformName]*internalPlatform
	deviceTree   map[adaptorName]map[platformName]map[deviceName]*internalDevice
}

var deviceTree _deviceTree

func errorNoSuchAdaptor(name string) error {
	return fmt.Errorf("There is no adaptor named %s on the adaptorTree", name)
}

func errorNoSuchPlatform(name string) error {
	return fmt.Errorf("There is no Platform named %s on the platformTree", name)
}

func errorNoSuchDevice(name string) error {
	return fmt.Errorf("There is no Device named %s on the deviceTree", name)
}

func (_this *_deviceTree) getAdaptorByName(targetAdaptor string) (res *internalAdaptor, err error) {
	_this.lock.RLock()
	res, ok := _this.adaptorTree[adaptorName(targetAdaptor)]
	_this.lock.RUnlock()
	if !ok {
		return nil, errorNoSuchAdaptor(targetAdaptor)
	}
	return res, nil
}

func (_this *_deviceTree) getPlatformByName(targetAdaptor string, targetPlatform string) (res *internalPlatform, err error) {
	_this.lock.RLock()
	res, ok := _this.platformTree[adaptorName(targetAdaptor)][platformName(targetPlatform)]
	_this.lock.RUnlock()
	if !ok {
		return nil, errorNoSuchPlatform(targetPlatform)
	}
	return res, nil
}
