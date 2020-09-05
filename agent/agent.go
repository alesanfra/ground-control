package agent

import "time"

// StartAgent starts the agent
func StartAgent(network string, httpPort uint) {
	devices := NewDeviceList()
	webServer := WebServer{devices: devices, port: httpPort}
	discoveryManager := DiscoveryManager{network: network, devices: devices}

	go webServer.Start()

	for {
		go discoveryManager.runDiscovery()
		<-time.After(30 * time.Second)
	}
}
