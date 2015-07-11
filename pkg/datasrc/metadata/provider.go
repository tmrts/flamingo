package metadata

import "fmt"

type URL string

func (u URL) WithVersion(v Version) string {
	return fmt.Sprintf(string(u), v)
}

type ProviderType string

const (
	GCE = "GoogleComputeEngine"
	EC2 = "ElasticComputeCloud"
)

type Source interface {
	MetaData(Version) (Interface, error)
}

type Provider struct {
	Name              string
	SupportedVersions map[Version]bool

	URL URL
}

func (p *Provider) MetaData(v Version) (Interface, error) {
	if _, ok := p.SupportedVersions[v]; ok != true {
		return nil, fmt.Errorf("metadata: version %v for %v is not supported", v, p.Name)
	}

	url := p.URL.WithVersion(v)
	metadata, err := request.Get(url, request.Header("Metadata-Flavor", "Google"))
	if err != nil {
		return nil, err
	}

	/*
	 *switch Provider {
	 *case "GoogleComputeEngine":
	 */
	var gceMetaData *GCEv1

	err = metadata.JSON(gceMetaData)
	if err != nil {
		return nil, err
	}

	return gceMetaData, nil
}
