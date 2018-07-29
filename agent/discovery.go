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

	nmap "github.com/lair-framework/go-nmap"
)

// DiscoveryManager manager
type DiscoveryManager struct {
	interval int
	network  string
	devices  *DeviceList
}

func (d *DiscoveryManager) runDiscovery() {
	ipv4Result := make(chan []nmap.Host, 1)
	ipv6Result := make(chan []nmap.Host, 1)
	go runDiscoveryIPv4(d.network, ipv4Result)
	go runDiscoveryIPv6(d.network, ipv6Result)
	ipv4Hosts := <-ipv4Result
	ipv6Hosts := <-ipv6Result
	d.devices.UpdateWithDiscoveryResult(ipv4Hosts, ipv6Hosts)
	d.logResults()
}

func runDiscoveryIPv4(network string, result chan []nmap.Host) {
	log.Print("Discovery IPv4 Start")

	binary, err := exec.LookPath("nmap")
	if err != nil {
		panic(err)
	}

	out, _ := exec.Command(binary, "-sn", "-oX", "-", network).Output()
	run, _ := nmap.Parse(out)
	result <- run.Hosts

	log.Print("Discovery IPv4 End")
}

func runDiscoveryIPv6(network string, result chan []nmap.Host) {
	log.Print("Discovery IPv6 Start")

	binary, err := exec.LookPath("nmap")
	if err != nil {
		panic(err)
	}

	out, _ := exec.Command(binary, "-6", "-sn", "-oX", "-", "--script=targets-ipv6-multicast-echo.nse", "--script-args", "'newtargets'").Output()
	run, _ := nmap.Parse(out)
	result <- run.Hosts

	log.Print("Discovery IPv6 End")
}

func (d *DiscoveryManager) logResults() {
	for _, host := range d.devices.hosts {
		log.Printf("%s: %s\n", host.Addresses[0].Addr, host.Status.State)
	}
}
