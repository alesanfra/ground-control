package connectivity

import (
	"context"
	"encoding/json"
	"log"
	"time"

	ndt7 "github.com/m-lab/ndt7-client-go"
)

// Ping the external world
func Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	client := ndt7.NewClient("ndt7-client-go-example", "0.1.0")
	ch, err := client.StartDownload(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for ev := range ch {
		log.Printf("%+v", ev)
	}
	ch, err = client.StartUpload(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for ev := range ch {
		log.Printf("%+v", ev)
	}

	a, _ := json.MarshalIndent(client.Results(), "", "    ")

	log.Print(string(a))
}
