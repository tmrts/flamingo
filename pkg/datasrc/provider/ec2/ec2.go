package ec2

import (
	"net"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/provider"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
	"github.com/tmrts/flamingo/pkg/request"
)

const (
	LatestSupportedMetadataVersion = "2009-04-04"

	MetadataURL provider.FormatURL = "http://169.254.169.254/%s/%s/%s"
)

type MetadataService struct {
	URL provider.FormatURL
}

func (s *MetadataService) fetchAttribute(attr string) (string, error) {
	return "", nil
}

func (s *MetadataService) FetchMetadata() (*metadata.Digest, error) {
	version := LatestSupportedMetadataVersion

	fetchMetadataAttribute := func(attr string) (string, error) {
		url := s.URL.Fill(version, "meta-data", attr)

		response, err := request.Get(url)
		if err != nil {
			return "", err
		}

		buf, err := response.Text()
		if err != nil {
			return "", err
		}

		return string(buf), nil
	}

	attributes := map[string]string{
		"hostname":                  "",
		"local-ipv4":                "",
		"public-ipv4":               "",
		"public-keys/0/openssh-key": "",
	}

	for attr := range attributes {
		value, err := fetchMetadataAttribute(attr)
		if err != nil {
			return nil, err
		}

		attributes[attr] = value
	}

	metadata := Metadata{
		Hostname:   attributes["hostname"],
		LocalIPv4:  net.ParseIP(attributes["local-ipv4"]),
		PublicIPv4: net.ParseIP(attributes["public-ipv4"]),
		PublicKeys: map[string]string{
			"openssh-key": attributes["public-keys/0/openssh-key"],
		},
	}

	digest := metadata.Digest()

	return &digest, nil
}

func (s *MetadataService) FetchUserdata() (userdata.Map, error) {
	// TODO(tmrts): Make sure if /user-data/ and /user-data paths are both valid.
	url := s.URL.Fill(LatestSupportedMetadataVersion, "user-data", "")

	response, err := request.Get(url[:len(url)-1])
	if err != nil {
		return nil, err
	}

	buf, err := response.Text()
	if err != nil {
		return nil, err
	}

	return userdata.Map{
		"user-data": string(buf),
	}, nil
}
