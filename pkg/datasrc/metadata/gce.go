package metadata

import "net"

type GCEv1 struct {
	ID          float64
	Image       string
	Hostname    string
	Description string
	CpuPlatform string
	MachineType string
	Zone        string

	MaintenanceEvent string
	Scheduling       struct {
		AutomaticRestart  string
		OnHostMaintenance string
	}

	VirtualClock struct {
		DriftToken string
	}

	NetworkInterfaces []struct {
		AccessConfigs []struct {
			ExternalIP net.IP
			Type       string
		}
		IP           net.IP
		Network      string
		ForwardedIPs []net.IP
	}

	Disks []struct {
		DeviceName string
		Index      int
		Mode       string
		Type       string
	}

	Attributes map[string]interface{}
	Tags       []string
}

func (metadata *GCEv1) Digest() Digest {
	interfaces := []Interface{}

	for _, intrfc := range metadata.NetworkInterfaces {
		i := Interface{
			NetworkName: intrfc.Network,
			PrivateIP:   intrfc.IP,
			PublicIPs:   []net.IP{},
		}
		for _, conf := range intrfc.AccessConfigs {
			i.PublicIPs = append(i.PublicIPs, conf.ExternalIP)
		}

		interfaces = append(interfaces, i)
	}

	return Digest{
		Hostname:          metadata.Hostname,
		NetworkInterfaces: interfaces,
	}
}
