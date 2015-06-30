package datasource

import (
	"fmt"
	"net"
	"strings"
)

type DataSource interface {
	FetchData() ([]byte, error)
}

type MetaData struct {
	PublicIPv4    net.IP
	PublicIPv6    net.IP
	PrivateIPv4   net.IP
	PrivateIPv6   net.IP
	HostName      string
	SSHPublicKeys map[string]string
	NetworkConfig []byte
}

const (
	DefaultAddress = "http://169.254.169.254/"
	apiVersion     = "2009-04-04/"
	userdataPath   = apiVersion + "user-data"
	metadataPath   = apiVersion + "meta-data"
)

func (md *MetaData) FetchData() (Metadata, error) {
	metadata := Metadata{}

	if keynames, err := md.fetchAttributes(fmt.Sprintf("%s/public-keys", md.MetadataUrl())); err == nil {
		keyIDs := make(map[string]string)
		for _, keyname := range keynames {
			tokens := strings.SplitN(keyname, "=", 2)
			if len(tokens) != 2 {
				return metadata, fmt.Errorf("malformed public key: %q", keyname)
			}
			keyIDs[tokens[1]] = tokens[0]
		}

		metadata.SSHPublicKeys = map[string]string{}
		for name, id := range keyIDs {
			sshkey, err := md.fetchAttribute(fmt.Sprintf("%s/public-keys/%s/openssh-key", md.MetadataUrl(), id))
			if err != nil {
				return metadata, err
			}
			metadata.SSHPublicKeys[name] = sshkey
			fmt.Printf("Found SSH key for %q\n", name)
		}
	} else if _, ok := err.(pkg.ErrNotFound); !ok {
		return metadata, err
	}

	if hostname, err := md.fetchAttribute(fmt.Sprintf("%s/hostname", md.MetadataUrl())); err == nil {
		metadata.Hostname = strings.Split(hostname, " ")[0]
	} else if _, ok := err.(pkg.ErrNotFound); !ok {
		return metadata, err
	}

	if localAddr, err := md.fetchAttribute(fmt.Sprintf("%s/local-ipv4", md.MetadataUrl())); err == nil {
		metadata.PrivateIPv4 = net.ParseIP(localAddr)
	} else if _, ok := err.(pkg.ErrNotFound); !ok {
		return metadata, err
	}

	if publicAddr, err := md.fetchAttribute(fmt.Sprintf("%s/public-ipv4", md.MetadataUrl())); err == nil {
		metadata.PublicIPv4 = net.ParseIP(publicAddr)
	} else if _, ok := err.(pkg.ErrNotFound); !ok {
		return metadata, err
	}

	return metadata, nil
}
