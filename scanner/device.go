package scanner

import (
	"net"
	"time"
)

type DeviceStatus string

const (
	Up   = DeviceStatus("up")
	Down = DeviceStatus("down")
	Lost = DeviceStatus("lost")
)

type Device struct {
	Ip       net.IP
	Mac      net.HardwareAddr
	Vendor   string
	LastSeen time.Time
	Status   DeviceStatus
}

type DeviceMap map[string]Device

func (dm DeviceMap) AsList() []Device {
	values := make([]Device, 0, len(dm))
	for _, value := range dm {
		values = append(values, value)
	}
	return values
}

func NewDeviceMap() DeviceMap {
	return make(map[string]Device)
}
