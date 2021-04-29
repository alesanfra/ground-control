package scanner

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ArpScanService struct {
	devices  DeviceMap
	lenience time.Duration
	interval time.Duration
}

func NewArpScanService(devices DeviceMap, lenience time.Duration, interval time.Duration) *ArpScanService {
	return &ArpScanService{
		devices:  devices,
		lenience: lenience,
		interval: interval,
	}
}

func (s *ArpScanService) Name() string {
	return "ARP scan"
}

func (s *ArpScanService) Run(ctx context.Context) error {
	// Get a list of all interfaces.
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	result := make(chan Device)
	var wg sync.WaitGroup

	for _, iface := range ifaces {
		// Start up a scan on each interface.
		go func(iface net.Interface) {
			wg.Add(1)
			defer wg.Done()
			if err := ArpScan(ctx, &iface, s.interval, result); err != nil {
				log.Printf("interface %v: %v", iface.Name, err)
			}
		}(iface)
	}

	go s.devices.SetDownAfter(s.lenience, ctx)
	go s.devices.AddDevices(result, ctx)

	wg.Wait()
	return err
}

// ArpScan scans an individual interface's local network for machines using ARP requests/replies.
// scan loops forever, sending packets out regularly.
func ArpScan(ctx context.Context, iface *net.Interface, interval time.Duration, result chan<- Device) error {
	addr, err := getIpAddress(iface)
	if err != nil {
		return err
	}

	log.Printf("Using network range %v for interface %v", addr, iface.Name)

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Start up a goroutine to read in packet data.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go readARP(ctx, handle, iface, result)

	for {
		select {
		case <-ctx.Done():
			log.Println("ArpScan: context canceled")
			return nil
		case <-time.After(interval):
			// Write our scan packets out to the handle.
			if err := writeARP(handle, iface, addr); err != nil {
				log.Printf("error writing packets on %v: %v", iface.Name, err)
			}
		}
	}
}

// We just look for IPv4 addresses, so try to find if the interface has one.
func getIpAddress(iface *net.Interface) (*net.IPNet, error) {

	if addresses, err := iface.Addrs(); err != nil {
		return nil, err
	} else {
		for _, a := range addresses {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {

					addr := &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}

					// Sanity-check that the interface has a good address.
					if addr.IP[0] == 127 {
						return nil, errors.New("skipping localhost")
					} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
						return nil, errors.New("mask means network is too large")
					}

					return addr, nil
				}
			}
		}
	}

	return nil, errors.New("no good IP network found")
}

// ReadARP watches a handle for incoming ARP responses we might care about, and prints them.
//
// ReadARP loops until context is closed
func readARP(ctx context.Context, handle *pcap.Handle, iface *net.Interface, result chan<- Device) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()

	for {
		var packet gopacket.Packet
		select {
		case <-ctx.Done():
			log.Println("ReadARP: context canceled")
			return
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			if arp.Operation != layers.ARPReply || bytes.Equal(iface.HardwareAddr, arp.SourceHwAddress) {
				// This is a packet I sent.
				continue
			}

			hwAddress := net.HardwareAddr(arp.SourceHwAddress).String()

			device := Device{
				Ip:       net.IP(arp.SourceProtAddress).String(),
				Mac:      hwAddress,
				Vendor:   GetVendor(hwAddress, true),
				LastSeen: time.Now().UTC(),
				Status:   Up,
			}
			// Note:  we might get some packets here that aren't responses to ones we've sent,
			// if for example someone else sends US an ARP request.  Doesn't much matter, though...
			// all information is good information :)
			result <- device
		}
	}
}

// WriteARP writes an ARP request for each address on our local network to the
// pcap handle.
func writeARP(handle *pcap.Handle, iface *net.Interface, addr *net.IPNet) error {
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(addr.IP),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// Send one packet for every address.
	for _, ip := range GetAllIpAddressesWithinNet(addr) {
		arp.DstProtAddress = ip
		if err := gopacket.SerializeLayers(buf, opts, &eth, &arp); err != nil {
			return err
		}
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

// GetAllIpAddressesWithinNet computes all IPv4 addresses contained in net.IPNet
func GetAllIpAddressesWithinNet(ipNetwork *net.IPNet) []net.IP {
	num := binary.BigEndian.Uint32(ipNetwork.IP)
	mask := binary.BigEndian.Uint32(ipNetwork.Mask)
	network := num & mask
	broadcast := network | ^mask

	// Pre-allocate result
	addresses := make([]net.IP, 0, broadcast-network-2)

	for network++; network < broadcast; network++ {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], network)
		addresses = append(addresses, buf[:])
	}

	return addresses
}
