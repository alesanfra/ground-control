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
