package gce

import (
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/request"
)

const (
	MetadataURL metadata.FormatURL = "http://metadata.google.internal/computeMetadata/%s/%s/?recursive=true"
)

type MetadataService struct {
	URL metadata.FormatURL
}

func (s *MetadataService) Metadata() (metadata.Interface, error) {
	// Currently only supported metadata version is "v1"
	version := "v1"
	/*
	 *if _, ok := s.SupportedVersions[v]; ok != true {
	 *    return nil, fmt.Errorf("metadata: version %v for %v is not supported", v, s.Name)
	 *}
	 */

	var m Metadata

	instanceUrl := s.URL.Fill(version, "instance")
	response, err := request.Get(instanceUrl, request.Header("Metadata-Flavor", "Google"))
	if err != nil {
		return nil, err
	}

	err = response.JSON(&m.Instance)
	if err != nil {
		return nil, err
	}

	projectUrl := s.URL.Fill(version, "project")
	response, err = request.Get(projectUrl, request.Header("Metadata-Flavor", "Google"))
	if err != nil {
		return nil, err
	}

	err = response.JSON(&m.Project)
	if err != nil {
		return nil, err
	}

	return &m, nil
}
