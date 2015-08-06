package gce

import (
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/provider"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
	"github.com/tmrts/flamingo/pkg/request"
)

const (
	LatestSupportedMetadataVersion = "v1"

	MetadataURL provider.FormatURL = "http://metadata.google.internal/computeMetadata/%s/%s/?recursive=true"
)

type MetadataService struct {
	URL provider.FormatURL
}

func (s *MetadataService) fetchMetadata() (*Metadata, error) {
	// currently only supported metadata version is "v1"
	version := "v1"
	/*
	 *if _, ok := s.supportedversions[v]; ok != true {
	 *    return nil, fmt.errorf("metadata: version %v for %v is not supported", v, s.name)
	 *}
	 */

	var m Metadata

	instanceurl := s.URL.Fill(version, "instance")
	response, err := request.Get(instanceurl, request.Header("Metadata-Flavor", "Google"))
	if err != nil {
		return nil, err
	}

	err = response.JSON(&m.Instance)
	if err != nil {
		return nil, err
	}

	projecturl := s.URL.Fill(version, "project")
	response, err = request.Get(projecturl, request.Header("Metadata-Flavor", "Google"))
	if err != nil {
		return nil, err
	}

	err = response.JSON(&m.Project)
	return &m, err
}

func (s *MetadataService) FetchMetadata() (*metadata.Digest, error) {
	m, err := s.fetchMetadata()
	if err != nil {
		return nil, err
	}

	digest := m.Digest()

	return &digest, nil
}

// FetchUserdata retrieves userdata files from Google Compute Engine metadata service
func (s *MetadataService) FetchUserdata() (userdata.Map, error) {
	m, err := s.fetchMetadata()
	return m.Instance.Attributes, err
}
