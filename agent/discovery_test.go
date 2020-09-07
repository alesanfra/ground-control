package agent

import (
	"bytes"
	"log"
	"testing"

	nmap "github.com/lair-framework/go-nmap"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLogResults(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	Convey("Given a discovery manager", t, func() {
		devices := NewDeviceList()
		devices.UpdateWithDiscoveryResult([]nmap.Host{
			{Comment: "test comment", Addresses: []nmap.Address{{Addr: "test_mac", AddrType: "mac"}}}},
			[]nmap.Host{})
		manager := DiscoveryManager{network: "127.0.0.1/32", devices: devices}

		Convey("It should be able to log results", func() {
			manager.logResults()
			So(len(buf.String()), ShouldBeGreaterThan, 0)
		})
	})
}

// func TestDiscovery(t *testing.T) {
// 	log.SetOutput(ioutil.Discard) // suppress log during test

// 	Convey("Given a discovery manager", t, func() {
// 		devices := &DeviceList{
// 			hosts: make(map[string]*nmap.Host),
// 		}
// 		manager := DiscoveryManager{network: "127.0.0.1/32", devices: devices}

// 		Convey("When runDiscovery is called", func() {
// 			res := manager.runDiscovery()

// 			Convey("Result should be a list of devices", func() {
// 				So(len(res), ShouldEqual, 0)
// 			})
// 		})
// 	})
// }
