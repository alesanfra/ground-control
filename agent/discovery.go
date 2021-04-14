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

func (d *DiscoveryManager) runDiscovery() []nmap.Host {
	ipv4Result := make(chan []nmap.Host, 1)
	ipv6Result := make(chan []nmap.Host, 1)
	go runDiscoveryIPv4(d.network, ipv4Result)
	go runDiscoveryIPv6(d.network, ipv6Result)
	ipv4Hosts := <-ipv4Result
	ipv6Hosts := <-ipv6Result
	d.devices.UpdateWithDiscoveryResult(ipv4Hosts, ipv6Hosts)
	d.logResults()

	ret := make([]nmap.Host, 0, len(d.devices.hosts))

	for _, value := range d.devices.hosts {
		ret = append(ret, *value)
	}
	return ret
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
