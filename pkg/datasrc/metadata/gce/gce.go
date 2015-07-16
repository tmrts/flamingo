package gce

import (
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/request"
)

const (
	projectMetadataURL  URLFormat = "http://metadata.google.internal/computeMetadata/%s/project/?recursive=true"
	instanceMetadataURL URLFormat = "http://metadata.google.internal/computeMetadata/%s/instance/?recursive=true"
)

type MetadataService struct {
	SupportedVersions map[Version]bool
}

func (s *MetadataService) Metadata(c request.Client) (metadata.Interface, error) {
	/*
	 *if _, ok := s.SupportedVersions[v]; ok != true {
	 *    return nil, fmt.Errorf("metadata: version %v for %v is not supported", v, s.Name)
	 *}
	 */

	var m *Metadata

	responce, err := request.Get(instanceMetadataURL.WithVersion("v1"), request.Header("Metadata-Flavor", "Google"))
	if err != nil {
		return nil, err
	}

	err = response.JSON(&m.Instance)
	if err != nil {
		return nil, err
	}

	pu := projectMetadataURL.WithVersion(v)
	responce, err := request.Get(projectMetadataURL.WithVersion("v1"), request.Header("Metadata-Flavor", "Google"))
	if err != nil {
		return nil, err
	}

	err = response.JSON(&m.Project)
	if err != nil {
		return nil, err
	}

	return m, nil
}
