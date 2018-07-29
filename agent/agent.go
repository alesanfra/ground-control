package agent

import "time"

// StartAgent starts the agent
func StartAgent(network string) {
	devices := NewDeviceList()
	webServer := WebServer{devices}
	discoveryManager := DiscoveryManager{network: network, devices: devices}

	go webServer.Start()

	for {
		go discoveryManager.runDiscovery()
		<-time.After(30 * time.Second)
	}
}
