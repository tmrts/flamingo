package metadata

var SupportedProviders = map[ProviderType]Source{
	GCE: Provider{
		Name:              "Google Compute Engine",
		SupportedVersions: []Version{"v1"},

		URL: "http://metadata.google.internal/computeMetadata/%s/instance/?recursive=true",
	},
	EC2: Provider{
		Name:              "Amazon Elastic Compute Cloud",
		SupportedVersions: []Version{"latest"},

		URL: "http://169.254.169.254/%s/meta-data/",
	},
}
