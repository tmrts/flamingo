package gce

import (
	"net"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

// Metadata represents the version "v1" compute meta-data provided
// by Google Compute Engine. Uninteresting fields are not exported.
type Metadata struct {
	Instance struct {
		iD          float64
		image       string
		Hostname    string
		description string
		cpuPlatform string
		machineType string
		zone        string

		maintenanceEvent string
		scheduling       struct {
			AutomaticRestart  string
			OnHostMaintenance string
		}

		virtualClock struct {
			DriftToken string
		}

		NetworkInterfaces []struct {
			IP           net.IP
			Network      string
			ForwardedIPs []net.IP

			AccessConfigs []struct {
				Type       string
				ExternalIP net.IP
			}
		}

		Disks []struct {
			Index int

			Type       string
			DeviceName string
			Mode       string
		}

		attributes map[string]interface{}

		tags []string
	}

	Project struct {
		ID         float64
		Attributes struct {
			SSHKeys []ssh.Key
		}
	}
}

// Digest extracts the important parts of meta-data and returns it.
func (md *Metadata) Digest() metadata.Digest {
	interfaces := []metadata.NetworkInterface{}
	for _, ifc := range md.Instance.NetworkInterfaces {
		i := metadata.NetworkInterface{
			NetworkName: ifc.Network,
			PrivateIP:   ifc.IP,
			PublicIPs:   []net.IP{},
		}

		for _, conf := range ifc.AccessConfigs {
			i.PublicIPs = append(i.PublicIPs, conf.ExternalIP)
		}

		interfaces = append(interfaces, i)
	}

	return metadata.Digest{
		Hostname: md.Instance.Hostname,
		SSHKeys:  md.Project.Attributes.SSHKeys,

		NetworkInterfaces: interfaces,
	}
}
