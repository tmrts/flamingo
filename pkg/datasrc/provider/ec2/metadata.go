package ec2

import (
	"net"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

type Metadata struct {
	instanceAction     string      `json:"instance-action"`
	instanceID         string      `json:"instance-id"`
	instanceType       string      `json:"instance-type"`
	kernelID           string      `json:"kernel-id"`
	localHostname      string      `json:"local-hostname"`
	publicHostname     string      `json:"public-hostname"`
	aMIID              string      `json:"ami-id"`
	aMILaunchIndex     string      `json:"ami-launch-index"`
	aMIManifestPath    string      `json:"ami-manifest-path"`
	blockDeviceMapping string      `json:"block-device-mapping"`
	ramDiskID          string      `json:"ramdisk-id"`
	reservationID      string      `json:"reservation-id"`
	securityGroups     []string    `json:"security-groups"`
	placement          interface{} `json:"placement"`

	Hostname   string            `json:"hostname"`
	LocalIPv4  net.IP            `json:"local-ipv4"`
	PublicIPv4 net.IP            `json:"public-ipv4"`
	PublicKeys map[string]string `json:"public-keys"`
}

// Digest extracts the important parts of meta-data and returns it.
func (m *Metadata) Digest() metadata.Digest {
	sshKeys := make(map[string][]ssh.Key)

	for _, publicKey := range m.PublicKeys {
		sshKeys["root"] = append(sshKeys["root"], ssh.Key(publicKey))
	}

	primaryNetworkInterface := metadata.NetworkInterface{
		PrivateIP: m.LocalIPv4,
		PublicIPs: []net.IP{m.PublicIPv4},
	}

	return metadata.Digest{
		Hostname: m.Hostname,
		SSHKeys:  sshKeys,

		NetworkInterfaces: []metadata.NetworkInterface{
			primaryNetworkInterface,
		},
	}
}
