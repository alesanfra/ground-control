// Ground Control: free and automated network scanner
// Copyright (C) 2018 Alessio Sanfratello
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
