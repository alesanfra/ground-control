package conf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type ArpScanConfig struct {
	Interval Duration
	Lenience Duration
}

type SpeedTestConfig struct {
	Interval Duration
}

type HttpConfig struct {
	Port int
}

type TelegramConfig struct {
	Token   string
	ChatId  int
	Timeout Duration
}

type Config struct {
	ArpScan   ArpScanConfig
	SpeedTest SpeedTestConfig
	Http      HttpConfig
	Telegram  TelegramConfig
}

// Read configuration from file
func Read(path string) (Config, error) {
	conf := Config{
		ArpScan: ArpScanConfig{
			Interval: Duration{15 * time.Second},
			Lenience: Duration{5 * time.Minute},
		},
		SpeedTest: SpeedTestConfig{
			Interval: Duration{2 * time.Hour},
		},
		Http: HttpConfig{
			Port: 5000,
		},
		Telegram: TelegramConfig{
			Timeout: Duration{1 * time.Hour},
		},
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return conf, err
	}

	err = json.Unmarshal(file, &conf)
	return conf, err
}
