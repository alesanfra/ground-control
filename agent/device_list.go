// Ground Control: free and automated network scanner
// Copyright (C) 2018 Alessio Sanfratello
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

// Init default constructor
func (d *DeviceList) Init() {
	d.hosts = make(map[string]*nmap.Host)
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
func (d *DeviceList) UpdateWithDiscoveryResult(hosts []nmap.Host) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	//mark all hosts as Down
	for mac := range d.hosts {
		d.hosts[mac].Status.State = "down"
	}

	// update device list with discovered hosts
	for i, host := range hosts {
		for _, address := range host.Addresses {
			if address.AddrType == "mac" {
				d.hosts[address.Addr] = &hosts[i]
				continue
			}
		}
	}
}
