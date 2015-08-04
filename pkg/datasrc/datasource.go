package datasrc

import (
	"errors"
	"log"
	"time"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/gce"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/openstack"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
)

type Provider interface {
	metadata.Provider
	userdata.Provider
}

// SupportedProviders returns the table of supported meta-data sources.
func SupportedProviders() map[string]Provider {
	return map[string]Provider{
		"GCE":       &gce.MetadataService{gce.MetadataURL},
		"OpenStack": &openstack.MetadataService{openstack.MetadataURL},
		//"OpenStack": configdrive.Mount,
		//"EC2": &ec2.MetadataService{openstack.MetadataURL},
		/*
		 *        Name: "Google Compute Engine",
		 *
		 *        URL: "http://metadata.google.internal/computeMetadata/%s/instance/?recursive=true",
		 *        "EC2": &DataSource{
		 *            Name: "Amazon Elastic Compute Cloud",
		 *            SupportedVersions: []Version{
		 *                "latest",
		 *            },
		 *
		 *            URL: "http://169.254.169.254/%s/meta-data/",
		 *        },
		 */
	}
}

var (
	ErrDatasourceRetrievalTimeout = errors.New("datasrc: timeout during data-source retrieval")
)

// isAvailable tries to fetch meta-data from the given datasource provider
// and returns the error if it encounters any.
func isAvailable(p Provider) bool {
	_, err := p.FetchMetadata()
	if err != nil {
		log.Print("datasrc.isAvailable:", err)
	}

	return err == nil
}

// FindProvider checks the given datasource providers and returns the available one.
func FindProvider(providers map[string]Provider, timeout time.Duration) (Provider, error) {
	providerChan := make(chan Provider)

	for _, p := range providers {
		go func() {
			if isAvailable(p) {
				providerChan <- p
			}
		}()
	}

	timeoutChan := time.NewTimer(timeout).C

	select {
	case p := <-providerChan:
		return p, nil
	case <-timeoutChan:
		return nil, ErrDatasourceRetrievalTimeout
	}
}
