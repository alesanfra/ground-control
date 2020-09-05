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

package agent

import (
	"fmt"
	"log"
	"net/http"
)

// WebServer rest api for the agent
type WebServer struct {
	devices *DeviceList
	port    uint
}

// Start starts http server
func (ws *WebServer) Start() {
	http.HandleFunc("/devices", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(ws.devices.AsJSON())
	})

	endpoint := fmt.Sprintf(":%d", ws.port)
	log.Printf("Start HTTP on %s\n", endpoint)
	http.ListenAndServe(endpoint, nil)

}
