/*
   Ground Control: free and automated network scanner
   Copyright (C) <2018>  <Alessio Sanfratello>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package agent

import (
	"bytes"
	"log"
	"os/exec"
	"time"

	"github.com/alesanfra/ground-control/cache"

	nmap "github.com/lair-framework/go-nmap"
)

const discoveryInterval = 30 * time.Second

// StartAgent starts the agent
func StartAgent(network string) {
	db := cache.MakeSolidCache("db")
	for {
		go networkDiscovery(network, db)
		<-time.After(discoveryInterval)
	}
}

func networkDiscovery(network string, db cache.SolidCache) {
	log.Print("Discovery Start")

	binary, lookErr := exec.LookPath("/usr/local/bin/nmap")
	if lookErr != nil {
		panic(lookErr)
	}

	cmd := exec.Command(binary, "nmap", "-sn", "-oX", "-", network)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	run, _ := nmap.Parse(out.Bytes())
	devices.UpdateWithDiscoveryResult(run.Hosts)
	go db.Put("roba", devices.hosts)

	for _, host := range devices.hosts {
		log.Printf("%s: %s\n", host.Addresses[0].Addr, host.Status.State)
	}
	log.Print("Discovery End")
}
