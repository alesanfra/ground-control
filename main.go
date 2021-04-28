package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alesanfra/ground-control/agent"
	"github.com/alesanfra/ground-control/scanner"
	"github.com/alesanfra/ground-control/web"
)

func main() {
	port := flag.Uint("p", 3000, "HTTP port")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	devices := scanner.NewDeviceMap()

	services := []agent.Service{
		scanner.NewArpScanService(devices, time.Minute, 10*time.Second),
		scanner.NewSpeedTestService(time.Minute),
		web.NewWebServer(devices, *port),
	}

	if err := agent.Run(ctx, services); err != nil {
		log.Fatalf("Error on agent run: %v", err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Print("Shutdown")
}
