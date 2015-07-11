package metadata

import (
	"net"

	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

// GCEv1 represented the v1 compute meta-data provided by google compute engine
// Uninteresting fields are unexported.
type GCEv1 struct {
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
func (metadata *GCEv1) Digest() Digest {
	interfaces := []NetworkInterface{}
	for _, ifc := range metadata.Instance.NetworkInterfaces {
		i := NetworkInterface{
			NetworkName: ifc.Network,
			PrivateIP:   ifc.IP,
			PublicIPs:   []net.IP{},
		}

		for _, conf := range ifc.AccessConfigs {
			i.PublicIPs = append(i.PublicIPs, conf.ExternalIP)
		}

		interfaces = append(interfaces, i)
	}

	return Digest{
		Hostname: metadata.Instance.Hostname,
		SSHKeys:  metadata.Project.Attributes.SSHKeys,

		NetworkInterfaces: interfaces,
	}
}
