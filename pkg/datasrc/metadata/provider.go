package metadata

import "fmt"

type URL string

func (u URL) WithVersion(v Version) string {
	return fmt.Sprintf(u, v)
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
	SupportedVersions []Version

	URL URL
}

func (p *Provider) MetaData(v Version) (Interface, error) {
	if _, ok := p.SupportedVersions[v]; ok != true {
		return nil, fmt.Errorf("metadata: version %v for %v is not supported", v, p.Name)
	}

	url := p.URL.WithVersion(v)
	metadata := request.Get(url, request.Header("Metadata-Flavor", "Google"))

	/*
	 *switch Provider {
	 *case "GoogleComputeEngine":
	 */
	var gceMetaData GCEv1

	metadata.JSON(&gceMetaData)

	return gceMetaData
}
