package scanner

import (
	"context"
	"log"
	"time"

	"github.com/m-lab/ndt7-client-go"
	"github.com/m-lab/ndt7-client-go/spec"
)

const (
	ClientName = "ground-control"
	Version    = "2"
)

type SpeedTestService struct {
	client   *ndt7.Client
	interval time.Duration
}

func NewSpeedTestService(interval time.Duration) *SpeedTestService {
	return &SpeedTestService{
		client:   ndt7.NewClient(ClientName, Version),
		interval: interval,
	}
}

func (s *SpeedTestService) Name() string {
	return "Speed Test"
}

func (s *SpeedTestService) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.interval):
			go func() {
				if results, err := SpeedTest(ctx, s.client); err != nil {
					log.Printf("Error on speed test: %v", err)
				} else {
					log.Printf("NDT: down %f Mbit/s, up %f Mbit/s, latency %f ms",
						results.Download, results.Upload, results.Latency)
				}
			}()
		}
	}
}

type SpeedTestResult struct {
	Download float64
	Upload   float64
	Latency  float64
}

// SpeedTest perform speed-test
func SpeedTest(ctx context.Context, client *ndt7.Client) (SpeedTestResult, error) {
	result := SpeedTestResult{Download: -1, Upload: -1, Latency: -1}

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	log.Println("NDT: start download")
	if download, err := client.StartDownload(ctx); err != nil {
		return result, err
	} else {
		for range download {
		}
	}

	log.Println("NDT: start upload")
	if upload, err := client.StartUpload(ctx); err != nil {
		return result, err
	} else {
		for range upload {
		}
	}

	rawResults := client.Results()

	// Download comes from the client-side Measurement during the download
	// test. DownloadRetrans and MinRTT come from the server-side Measurement,
	// if it includes a TCPInfo object.
	if dl, ok := rawResults[spec.TestDownload]; ok {
		if dl.Client.AppInfo != nil && dl.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(dl.Client.AppInfo.ElapsedTime) / 1e06
			result.Download = (8.0 * float64(dl.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)
		}
		if dl.Server.TCPInfo != nil {
			result.Latency = float64(dl.Server.TCPInfo.MinRTT) / 1000
		}
	}

	// Upload comes from the client-side Measurement during the upload test.
	if ul, ok := rawResults[spec.TestUpload]; ok {
		if ul.Client.AppInfo != nil && ul.Client.AppInfo.ElapsedTime > 0 {
			elapsed := float64(ul.Client.AppInfo.ElapsedTime) / 1e06
			result.Upload = (8.0 * float64(ul.Client.AppInfo.NumBytes)) / elapsed / (1000.0 * 1000.0)
		}
	}

	return result, nil
}
