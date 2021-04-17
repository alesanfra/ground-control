package scanner

import (
	"context"
	"log"
	"time"

	"github.com/m-lab/ndt7-client-go"
	"github.com/m-lab/ndt7-client-go/spec"
)

type SpeedTestResult struct {
	Download float64
	Upload   float64
	Latency  float64
}

// Ping the external world
func SpeedTest(ctx context.Context, client *ndt7.Client) SpeedTestResult {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	log.Println("NDT: start download")
	download, err := client.StartDownload(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for range download {
	}

	log.Println("NDT: start upload")
	upload, err := client.StartUpload(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for range upload {
	}

	results := makeSummary(client.Results())
	log.Printf("NDT: results: %v", results)
	return results
}

func makeSummary(results map[spec.TestKind]*ndt7.LatestMeasurements) SpeedTestResult {
	r := SpeedTestResult{Download: -1, Upload: -1, Latency: -1}

	// Download comes from the client-side Measurement during the download
	// test. DownloadRetrans and MinRTT come from the server-side Measurement,
	// if it includes a TCPInfo object.
	if dl, ok := results[spec.TestDownload]; ok {
		if dl.Client.AppInfo != nil && dl.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(dl.Client.AppInfo.ElapsedTime) / 1e06
			r.Download = (8.0 * float64(dl.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)
		}
		if dl.Server.TCPInfo != nil {
			r.Latency = float64(dl.Server.TCPInfo.MinRTT) / 1000
		}
	}

	// Upload comes from the client-side Measurement during the upload test.
	if ul, ok := results[spec.TestUpload]; ok {
		if ul.Client.AppInfo != nil && ul.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(ul.Client.AppInfo.ElapsedTime) / 1e06
			r.Upload = (8.0 * float64(ul.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)
		}
	}

	return r
}
