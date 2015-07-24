package ec2

import (
	"net"

	"github.com/tmrts/flamingo/pkg/datasrc/provider"
)

const (
	// currently only supported metadata version is "2012-08-10"
	LatestSupportedMetadataVersion string = "2009-04-04"

	MetadataURL provider.FormatURL = "http://169.254.169.254/meta-data/%s"
)

type Metadata struct {
	InstanceAction     string      `json:"instance-action"`
	InstanceID         string      `json:"instance-id"`
	InstanceType       string      `json:"instance-type"`
	KernelID           string      `json:"kernel-id"`
	Hostname           string      `json:"hostname"`
	LocalHostname      string      `json:"local-hostname"`
	PublicHostname     string      `json:"public-hostname"`
	AMIID              string      `json:"ami-id"`
	AMILaunchIndex     string      `json:"ami-launch-index"`
	AMIManifestPath    string      `json:"ami-manifest-path"`
	BlockDeviceMapping string      `json:"block-device-mapping"`
	PublicKeys         string      `json:"public-keys"`
	RamDiskID          string      `json:"ramdisk-id"`
	ReservationID      string      `json:"reservation-id"`
	LocalIPv4          net.IP      `json:"local-ipv4"`
	PublicIPv4         net.IP      `json:"public-ipv4"`
	SecurityGroups     []string    `json:"security-groups"`
	Placement          interface{} `json:"placement"`
}
