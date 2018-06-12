package main

import (
	"testing"

	nmap "github.com/lair-framework/go-nmap"
)

var expectedJSON string = `{"test_mac":{"starttime":"0001-01-01 00:00:00 +0000 UTC","endtime":"0001-01-01 00:00:00 +0000 UTC","comment":"test comment","status":{"state":"","reason":"","reason_ttl":0},"addresses":[{"addr":"test_mac","addrtype":"mac","vendor":""}],"hostnames":null,"smurfs":null,"ports":null,"extraports":null,"os":{"portsused":null,"osmatches":null,"osfingerprints":null},"distance":{"value":0},"uptime":{"seconds":0,"lastboot":""},"tcpsequence":{"index":0,"difficulty":"","vaules":""},"ipidsequence":{"class":"","values":""},"tcptssequence":{"class":"","values":""},"hostscripts":null,"trace":{"proto":"","port":0,"hops":null},"times":{"srtt":"","rttv":"","to":""}}}`

func TestAsJSONEmpty(t *testing.T) {
	var testDevices DeviceList
	testDevices.Init()
	j := testDevices.AsJSON()
	s := string(j)
	if s != "{}" {
		t.Error("Expected", "{}", "Got", s)
	}
}

func TestAsJSON(t *testing.T) {
	var testDevices DeviceList
	testDevices.Init()
	testDevices.UpdateWithDiscoveryResult([]nmap.Host{
		{Comment: "test comment", Addresses: []nmap.Address{{Addr: "test_mac", AddrType: "mac"}}}})
	j := testDevices.AsJSON()
	s := string(j)
	if s != expectedJSON {
		t.Error("Expected", expectedJSON, "Got", s)
	}
}
