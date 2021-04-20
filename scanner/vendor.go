package scanner

import (
	"log"

	"github.com/umahmood/macvendors"
)

const DefaultVendor = "Unknown"

type VendorFinder struct {
	vendorApi *macvendors.API
	cache     map[string]string
}

func NewVendorFinder() *VendorFinder {
	return &VendorFinder{
		vendorApi: macvendors.New(),
		cache:     make(map[string]string),
	}
}

func (f *VendorFinder) Find(hwAddress string) string {
	name, ok := f.cache[hwAddress]

	if !ok {
		var err error
		log.Printf("Vendor cache miss: %s not found, invoking 3rd party api", hwAddress)
		if name, err = f.vendorApi.Name(hwAddress); err != nil {
			log.Printf("Name for %s address not found, fallback to %s", hwAddress, DefaultVendor)
			name = DefaultVendor
		} else {
			f.cache[hwAddress] = name
		}
	}

	return name
}
