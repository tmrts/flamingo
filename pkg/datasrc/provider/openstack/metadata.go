package openstack

import (
	"strings"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

type file struct {
	ContentPath string `json:"content_path"`
	Path        string `json:"content_path"`
}

type Metadata struct {
	Name             string            `json:"name"`
	UUID             string            `json:"uuid"`
	Hostname         string            `json:"hostname"`
	ProjectID        string            `json:"project_id"`
	LaunchIndex      int               `json:"launch_index"`
	AvailabilityZone string            `json:"availability_zone"`
	PublicKeys       map[string]string `json:"public_keys"`
	Meta             map[string]string `json:"meta"`
	Files            []file            `json:"files"`
}

// Digest extracts the important parts of meta-data and returns it.
func (md *Metadata) Digest() metadata.Digest {
	sshKeys := make(map[string][]ssh.Key)

	for usr, rawKeys := range md.PublicKeys {
		keys := strings.Split(rawKeys, "\n")

		for _, key := range keys {
			if key != "" {
				sshKeys[usr] = append(sshKeys[usr], ssh.Key(key))
			}
		}
	}

	return metadata.Digest{
		Hostname: md.Hostname,
		SSHKeys:  sshKeys,
	}
}
