package metadata

import (
	"time"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/metadata/gce"
	"github.com/tmrts/flamingo/pkg/request"
)

// SupportedProviders is the table of supported meta-data sources.
var SupportedProviders = map[string]metadata.Provider{
	"GCE":       gce.MetadataService,
	"Openstack": configdrive.Mount,
	/*
	 *        Name: "Google Compute Engine",
	 *
	 *        URL: "http://metadata.google.internal/computeMetadata/%s/instance/?recursive=true",
	 */
	/*
	 *    "EC2": &Provider{
	 *        Name: "Amazon Elastic Compute Cloud",
	 *        SupportedVersions: []Version{
	 *            "latest",
	 *        },
	 *
	 *        URL: "http://169.254.169.254/%s/meta-data/",
	 *    },
	 */
}

// Get queries metadata services and returns the available metadata digest
// Since there will be only one metadata service operational (or multiple ones returning the same),
// it queries each known service and returns at timeout or a successful response.
func Get(timeout time.Duration) (digest metadata.Digest) {
	if configdrive.Available() {
		return configdrive.Metadata()
	}

	metadataChan := make(chan Digest, len(SupportedProviders))
	defer close(metadataChan)

	go func(mdchan chan Digest) {
		for _, p := range SupportedProviders {
			go func(c chan Digest) {
				md, err := p.Metadata(request.DefaultClient)
				if err != nil {
					return
				}

				c <- md.Digest()
			}(mdchan)
		}

		time.Sleep(timeout)
		c <- nil
		close(c)
	}(metadataChan)

	return <-metadataChan
}
