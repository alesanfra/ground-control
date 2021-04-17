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

type DeviceRegistry struct {
	Hosts map[string]Device
}

func (dr *DeviceRegistry) AsList() []Device {
	values := make([]Device, 0, len(dr.Hosts))
	for _, value := range dr.Hosts {
		values = append(values, value)
	}
	return values
}
