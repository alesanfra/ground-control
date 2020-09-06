package agent

// import (
// 	"io/ioutil"
// 	"log"
// 	"testing"

// 	nmap "github.com/lair-framework/go-nmap"
// 	. "github.com/smartystreets/goconvey/convey"
// )

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
