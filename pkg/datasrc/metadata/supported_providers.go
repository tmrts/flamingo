package metadata

var SupportedProviders = map[ProviderType]Source{
	GCE: &Provider{
		Name:              "Google Compute Engine",
		SupportedVersions: map[Version]bool{"v1": true},

		URL: "http://metadata.google.internal/computeMetadata/%s/instance/?recursive=true",
	},
	EC2: &Provider{
		Name:              "Amazon Elastic Compute Cloud",
		SupportedVersions: map[Version]bool{"latest": true},

		URL: "http://169.254.169.254/%s/meta-data/",
	},
}
