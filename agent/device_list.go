package agent

import (
	"encoding/json"
	"log"
	"sync"

	nmap "github.com/lair-framework/go-nmap"
)

// DeviceList protected list of hosts
type DeviceList struct {
	hosts map[string]*nmap.Host
	mutex sync.RWMutex
}

// NewDeviceList constructor
func NewDeviceList() *DeviceList {
	return &DeviceList{
		hosts: make(map[string]*nmap.Host),
	}
}

// AsJSON returns device list as json buffer
func (d *DeviceList) AsJSON() []byte {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	b, err := json.Marshal(d.hosts)
	if err != nil {
		log.Print("error:", err)
	}

	return b
}

// UpdateWithDiscoveryResult add discovery result
func (d *DeviceList) UpdateWithDiscoveryResult(ipv4Hosts []nmap.Host, ipv6Hosts []nmap.Host) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	//mark all hosts as Down
	for mac := range d.hosts {
		d.hosts[mac].Status.State = "down"
	}

	// update device list with discovered hosts
	d.addDiscoveredDevices(ipv4Hosts)
	d.addDiscoveredDevices(ipv6Hosts)
}

func (d *DeviceList) addDiscoveredDevices(hosts []nmap.Host) {
	for i, host := range hosts {
		for j, address := range host.Addresses {
			if address.AddrType == "mac" {
				d.hosts[hosts[i].Addresses[j].Addr] = &hosts[i]
				continue
			}
		}
	}
}
