package scanner

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const DefaultVendor = "Unknown"
const VendorNameEndPoint = "https://macvendors.co/api/vendorname/%s"

type VendorFinder struct {
	cache map[string]string
}

func NewVendorFinder() *VendorFinder {
	return &VendorFinder{
		cache: make(map[string]string),
	}
}

func (f *VendorFinder) Find(hwAddress string) string {
	name, ok := f.cache[hwAddress]
	if ok {
		return name
	}

	log.Printf("Vendor cache miss: %s not found, invoking 3rd party api", hwAddress)
	name, err := getVendor(hwAddress)
	if err != nil {
		log.Printf("Name for %s address not found, fallback to %s, error %v", hwAddress, DefaultVendor, err)
		return DefaultVendor
	}

	f.cache[hwAddress] = name
	return name
}

func getVendor(addr string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(VendorNameEndPoint, addr))
	if err != nil || resp.StatusCode >= 300 {
		return "", errors.New("error on vendor api call")
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("bad response from vendor api")
	}

	return string(body), nil
}
