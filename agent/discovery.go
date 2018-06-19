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
	"log"
	"os/exec"
	"time"

	nmap "github.com/lair-framework/go-nmap"
)

const discoveryInterval = 30 * time.Second

// StartAgent starts the agent
func StartAgent(network string) {
	devices := NewDeviceList()

	go StartHTTPServer(devices)

	for {
		go networkDiscovery(network, devices)
		go networkDiscoveryIPv6(devices)
		<-time.After(discoveryInterval)
	}
}

func networkDiscovery(network string, devices *DeviceList) {
	log.Print("Discovery Start")

	binary, err := exec.LookPath("nmap")
	if err != nil {
		panic(err)
	}

	out, _ := exec.Command(binary, "-sn", "-oX", "-", network).Output()
	run, _ := nmap.Parse(out)

	devices.UpdateWithDiscoveryResult(run.Hosts)
	logResults(devices)

	log.Print("Discovery End")
}

func networkDiscoveryIPv6(devices *DeviceList) {
	log.Print("Discovery Start")

	binary, err := exec.LookPath("nmap")
	if err != nil {
		panic(err)
	}

	out, _ := exec.Command(binary, "-6", "-oX", "-", "--script=targets-ipv6-multicast-echo.nse", "--script-args", "'newtargets'", "-sL").Output()
	run, _ := nmap.Parse(out)

	// devices.UpdateWithDiscoveryResult(run.Hosts)
	// logResults(devices)
	log.Print(run.Hosts)

	log.Print("Discovery End")

	//
}

func logResults(devices *DeviceList) {
	for _, host := range devices.hosts {
		log.Printf("%s: %s\n", host.Addresses[0].Addr, host.Status.State)
	}
}
