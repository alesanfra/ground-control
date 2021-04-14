package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alesanfra/ground-control/agent"
	"github.com/alesanfra/ground-control/connectivity"
)

func main() {
	connectivity.Ping()

	network := flag.String("n", "", "Network to be scanned in the form 192.168.1.0/24")
	port := flag.Uint("p", 3000, "HTTP port")
	flag.Parse()

	if len(*network) == 0 {
		log.Fatal("You must specify the network")
	}

	log.Printf("Start Network Discovery on %s\n", *network)

	go agent.StartAgent(*network, *port)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Print("Shutdown")

}
