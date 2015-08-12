// Package datasrc contains the interface required for being a data source
// provider and the implementations for finding an available data source.
package datasrc

import (
	"errors"
	"time"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/configdrive"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/ec2"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/gce"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/openstack"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
	"github.com/tmrts/flamingo/pkg/flog"
	"github.com/tmrts/flamingo/pkg/sys"
)

var (
	ErrDatasourceRetrievalTimeout = errors.New("datasrc: timeout during data-source retrieval")
)

// Provider is the interface that represents a data source provider that
// can provide meta-data and user-data information on a cloud service.
type Provider interface {
	metadata.Provider
	userdata.Provider
}

// SupportedProviders returns the table of supported data source providers.
func SupportedProviders() map[string]Provider {
	providers := map[string]Provider{
		"GCE":       &gce.MetadataService{gce.MetadataURL},
		"EC2":       &ec2.MetadataService{ec2.MetadataURL},
		"OpenStack": &openstack.MetadataService{openstack.MetadataURL},
	}

	// OpenStack Config-Drive overrides OpenStack meta-data service as the data provider.
	configdrivePath, err := configdrive.FindMountTarget(sys.DefaultExecutor)
	if err == nil {
		// TODO(tmrts): Decide whether config-drive should merge or
		//				override OpenStack meta-data service
		providers["OpenStack"] = &configdrive.Mount{configdrivePath}
	}

	return providers
}

// isAvailable tries to fetch meta-data from the given data source provider
// and returns false if it encounters any errors during the process.
func isAvailable(p Provider) bool {
	_, err := p.FetchMetadata()
	if err != nil {
		flog.Debug("Provider isn't available",
			flog.Fields{
				Event: "datasrc.isAvailable",
			},
			flog.Details{
				"provider": p,
			},
		)
	}

	return err == nil
}

// FindProvider checks the given datasource providers, if it finds an available
// data source before the specified duration, it returns the provider, else it
// returns an ErrDatasourceRetrievalTimeout error.
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
