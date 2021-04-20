package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alesanfra/ground-control/scanner"
	"github.com/alesanfra/ground-control/web"
	"github.com/m-lab/ndt7-client-go"
)

const ClientName = "ground-control"
const Version = "2"

func main() {
	port := flag.Uint("p", 3000, "HTTP port")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devices := scanner.NewDeviceMap()

	// start web server
	go web.NewWebServer(devices, *port).Start()
	go startArpScanner(ctx, devices)

	_ = ndt7.NewClient(ClientName, Version)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Print("Shutdown")

}

// StartAgent starts the agent
func startArpScanner(ctx context.Context, devices scanner.DeviceMap) {
	vendor := scanner.NewVendorFinder()

	// Get a list of all interfaces.
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	result := make(chan scanner.Device)

	for _, iface := range ifaces {
		// Start up a scan on each interface.
		go func(iface net.Interface) {
			if err := scanner.ArpScan(&iface, 10*time.Second, vendor, result, ctx); err != nil {
				log.Printf("interface %v: %v", iface.Name, err)
			}
		}(iface)
	}

	for device := range result {
		log.Printf("IP %v is at %v (%s)", device.Ip, device.Mac, device.Vendor)
		devices[device.Mac.String()] = device
	}
}
