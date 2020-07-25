package agent

import (
	"testing"

	nmap "github.com/lair-framework/go-nmap"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAsJSON(t *testing.T) {
	Convey("Given a device list", t, func() {
		testDevices := NewDeviceList()

		Convey("When the list is empty", func() {
			j := string(testDevices.AsJSON())

			Convey("The resulting json should be an empty object", func() {
				So(j, ShouldEqual, "{}")
			})
		})

		Convey("When the list contains devices", func() {
			testDevices.UpdateWithDiscoveryResult([]nmap.Host{
				{Comment: "test comment", Addresses: []nmap.Address{{Addr: "test_mac", AddrType: "mac"}}}},
				[]nmap.Host{})
			j := string(testDevices.AsJSON())

			Convey("The resulting json should contain devices as well", func() {
				expectedJSON := `{"test_mac":{"starttime":-62135596800,"endtime":-62135596800,"comment":"test comment","status":{"state":"","reason":"","reason_ttl":0},"addresses":[{"addr":"test_mac","addrtype":"mac","vendor":""}],"hostnames":null,"smurfs":null,"ports":null,"extraports":null,"os":{"portsused":null,"osmatches":null,"osfingerprints":null},"distance":{"value":0},"uptime":{"seconds":0,"lastboot":""},"tcpsequence":{"index":0,"difficulty":"","vaules":""},"ipidsequence":{"class":"","values":""},"tcptssequence":{"class":"","values":""},"hostscripts":null,"trace":{"proto":"","port":0,"hops":null},"times":{"srtt":"","rttv":"","to":""}}}`
				So(j, ShouldEqual, expectedJSON)
			})
		})
	})
}
