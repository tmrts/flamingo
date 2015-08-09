package main

import (
	"flag"
	"time"

	"github.com/tmrts/flamingo/pkg/datasrc"
	"github.com/tmrts/flamingo/pkg/distro"
	"github.com/tmrts/flamingo/pkg/flog"
	"github.com/tmrts/flamingo/pkg/sys"
)

var flags struct {
	cloudConfig string
	configDrive string
	metadata    string
}

func init() {
	flag.StringVar(&flags.cloudConfig, "cloud-config", "", "user-data configuration file")
	flag.StringVar(&flags.configDrive, "config-drive", "", "config drive mount path")
	flag.StringVar(&flags.metadata, "meta-data", "", "meta-data file")
}

func InitializeContextualization() error {
	// TODO(tmrts): Add plug-in hooks
	flog.Info("Initializing contextualization",
		flog.Fields{
			Event: "main.InitializeContextualization",
		},
	)
	return nil
}

func FinalizeContextualization() error {
	// TODO(tmrts): Add plug-in hooks
	flog.Info("Finalizing contextualization",
		flog.Fields{
			Event: "main.FinalizeContextualization",
		},
	)
	return nil
}

func HasRootPrivileges() (bool, error) {
	// TODO(tmrts): Move privilege check to sys.nss package
	ent, err := sys.Execute("getent", "gshadow", "root")
	if err != nil {
		return false, err
	}

	return ent != "", nil
}

func main() {
	// TODO(tmrts): Add command-line flags

	// TODO(tmrts): Build a logger hierarchy
	hasRoot, err := HasRootPrivileges()
	if err != nil {
		flog.Fatal("Failed checking user privileges",
			flog.Fields{
				Event: "main.HasRootPrivileges",
				Error: err,
			},
		)
	}

	if !hasRoot {
		flog.Fatal("current user doesn't have root privileges")
	}

	centOS := distro.CentOS(sys.DefaultExecutor)

	if err := InitializeContextualization(); err != nil {
		flog.Fatal("Failed to start contextualization",
			flog.Fields{
				Event: "main.InitializeContextualization",
				Error: err,
			},
		)
	}

	providers := datasrc.SupportedProviders()

	p, err := datasrc.FindProvider(providers, 5*time.Second)
	if err != nil {
		flog.Fatal("Failed to start contextualization",
			flog.Fields{
				Event: "datasrc.FindProvider",
				Error: err,
			},
		)
	}

	m, err := p.FetchMetadata()
	if err != nil {
		flog.Fatal("Failed to fetch meta-data from provider",
			flog.Fields{
				Event: "datasrc.Provider.FetchMetadata",
				Error: err,
			},
			flog.Details{
				"provider": p,
			},
		)
	}

	u, err := p.FetchUserdata()
	if err != nil {
		flog.Fatal("Failed to fetch user-data from provider",
			flog.Fields{
				Event: "datasrc.Provider.FetchUserdata",
				Error: err,
			},
			flog.Details{
				"provider": p,
			},
		)
	}

	if err := centOS.ConsumeMetadata(m); err != nil {
		flog.Fatal("Failed to consume meta-data",
			flog.Fields{
				Event: "distro.Implementation.ConsumeMetadata",
				Error: err,
			},
			flog.Details{
				"distribution": "CentOS",
			},
		)
	}

	if err := centOS.ConsumeUserdata(u); err != nil {
		flog.Fatal("Failed to consume user-data",
			flog.Fields{
				Event: "distro.Implementation.ConsumeMetadata",
				Error: err,
			},
			flog.Details{
				"distribution": "CentOS",
			},
		)
	}

	if err := FinalizeContextualization(); err != nil {
		flog.Fatal("Failed to finalize contextualization",
			flog.Fields{
				Event: "main.FinalizeContextualization",
				Error: err,
			},
		)
	}
}
