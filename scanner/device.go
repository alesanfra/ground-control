package scanner

import (
	"context"
	"log"
	"time"
)

type DeviceStatus string

const (
	Up   = DeviceStatus("up")
	Down = DeviceStatus("down")
)

type Device struct {
	Ip       string
	Mac      string
	Vendor   string
	LastSeen time.Time
	Status   DeviceStatus
}

type DeviceMap map[string]Device

func NewDeviceMap() DeviceMap {
	return make(map[string]Device)
}

func (dm DeviceMap) AsList() []Device {
	values := make([]Device, 0, len(dm))
	for _, value := range dm {
		values = append(values, value)
	}
	return values
}

func (dm DeviceMap) SetDownAfter(leniency time.Duration, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(10 * time.Second):
			now := time.Now().UTC()
			for _, device := range dm {
				if now.Sub(device.LastSeen) > leniency {
					device.Status = Down
					log.Printf("Set device %s (%s) to down, last seen %v", device.Mac, device.Vendor, device.LastSeen)
				}
			}
		}
	}

}
