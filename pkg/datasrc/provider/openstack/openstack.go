package openstack

import (
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/provider"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
	"github.com/tmrts/flamingo/pkg/request"
)

const (
	// currently only supported metadata version
	LatestSupportedMetadataVersion string = "2012-08-10"

	MetadataURL provider.FormatURL = "http://169.254.169.254/openstack/%s/%s"
)

type MetadataService struct {
	URL provider.FormatURL
}

// FetchMetadata retrieves meta_data.json from OpenStack Metadata
// service and parses it.
func (s *MetadataService) FetchMetadata() (*metadata.Digest, error) {
	var m Metadata

	url := s.URL.Fill(LatestSupportedMetadataVersion, "meta_data.json")
	response, err := request.Get(url)
	if err != nil {
		return nil, err
	}

	err = response.JSON(&m)
	if err != nil {
		return nil, err
	}

	digest := m.Digest()

	return &digest, nil
}

// FetchMetadata retrieves meta_data.json from OpenStack Metadata
// service and parses it.
func (s *MetadataService) FetchUserdata() (userdata.Map, error) {
	url := s.URL.Fill(LatestSupportedMetadataVersion, "user_data")
	response, err := request.Get(url)
	if err != nil {
		return nil, err
	}

	buf, err := response.Text()
	if err != nil {
		return nil, err
	}

	u := userdata.Map{
		"user-data": string(buf),
	}

	return u, nil
}
