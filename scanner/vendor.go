package scanner

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	DefaultVendor      = "Unknown"
	VendorNameEndPoint = "https://macvendors.co/api/vendorname/%s"
)

var vendorCache = make(map[string]string)

// GetVendor tries to retrieve the vendor of an interface from the mac address
func GetVendor(hwAddress string, cached bool) string {
	if cached {
		name, ok := vendorCache[hwAddress]
		if ok {
			return name
		} else {
			log.Printf("Vendor cache miss: %s not found, invoking 3rd party api", hwAddress)
		}
	}

	name, err := GetVendorFromApi(hwAddress)
	if err != nil {
		log.Printf("Name for %s address not found, fallback to %s, error %v", hwAddress, DefaultVendor, err)
		return DefaultVendor
	}

	vendorCache[hwAddress] = name
	return name
}

func GetVendorFromApi(addr string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(VendorNameEndPoint, addr))
	if err != nil || resp.StatusCode >= 300 {
		return DefaultVendor, errors.New("error on vendor api call")
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DefaultVendor, errors.New("bad response from vendor api")
	}

	return string(body), nil
}
