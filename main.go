package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/alesanfra/ground-control/agent"
	"github.com/alesanfra/ground-control/conf"
	"github.com/alesanfra/ground-control/notification"
	"github.com/alesanfra/ground-control/scanner"
	"github.com/alesanfra/ground-control/web"
)

func main() {
	port := flag.Uint("p", 3000, "HTTP port")
	configPath := flag.String("c", "conf.json", "Telegram configuration file")
	flag.Parse()

	configuration, err := conf.Read(*configPath)
	if err != nil {
		log.Fatalf("Unable to read configuration from %s", *configPath)
	}

	log.Printf("Read configuration %v", configuration)

	devices := scanner.NewDeviceMap()

	services := []agent.Service{
		scanner.NewArpScanService(devices, configuration.ArpScan.Lenience.Duration, configuration.ArpScan.Interval.Duration),
		scanner.NewSpeedTestService(time.Hour),
		web.NewWebService(devices, *port),
		notification.NewNotificationService(devices, configuration.Telegram),
	}

	if err := agent.Run(context.Background(), services); err != nil {
		log.Fatalf("Error on agent run: %v", err)
	}

	log.Print("Shutdown")
}
