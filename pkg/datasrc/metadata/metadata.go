package metadata

import (
	"fmt"
	"net"

	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

// Version is represents the version of a meta-data.
type Version string

// Interface represents the meta-data object returned
// by data sources that can be summarized into a meta-data digest.
type Interface interface {
	Digest() Digest
}

type Provider interface {
	FetchMetadata() (*Digest, error)
}

// Digest contains the parts of a meta-data object that are
// used by Flamingo for contextualization of the instance.
type Digest struct {
	Hostname string

	NetworkInterfaces []NetworkInterface

	SSHKeys map[string][]ssh.Key
}

func (d *Digest) String() string {
	return fmt.Sprintf("\nHostname: %v\nInterfaces: %v\n", d.Hostname, d.NetworkInterfaces)
}

func (d *Digest) PrimaryNetworkInterface() *NetworkInterface {
	return &d.NetworkInterfaces[0]
}

type NetworkInterface struct {
	NetworkName string

	PrivateIP net.IP

	PublicIPs []net.IP
}

func (i NetworkInterface) String() string {
	return fmt.Sprintf("\nNetworkName: %v\nPrivateIP: %v\nPublicIPs: %v\n", i.NetworkName, i.PrivateIP, i.PublicIPs)
}

type Disk struct {
	Mode       string
	Type       string
	DeviceName string
}
