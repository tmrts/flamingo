package metadata

import (
	"fmt"
	"net"
)

type Version string

type Interface interface {
	Digest() Digest
}

// Digest contains the summarized parts of a meta-data digest for consumption
type Digest struct {
	Hostname          string
	NetworkInterfaces []Interface
}

func (d Digest) String() string {
	return fmt.Sprintf("\nHostname: %v\nInterfaces: %v\n", d.Hostname, d.NetworkInterfaces)
}

type Interface struct {
	NetworkName string

	PrivateIP net.IP

	PublicIPs []net.IP
}

func (i Interface) String() string {
	return fmt.Sprintf("\nNetworkName: %v\nPrivateIP: %v\nPublicIPs: %v\n", i.NetworkName, i.PrivateIP, i.PublicIPs)
}

type Disk struct {
	Mode       string
	Type       string
	DeviceName string
}
