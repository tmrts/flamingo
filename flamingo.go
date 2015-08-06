package main

import (
	"flag"
	"time"

	"github.com/tmrts/flamingo/pkg/datasrc"
	"github.com/tmrts/flamingo/pkg/distro"
	"github.com/tmrts/flamingo/pkg/flog"
	"github.com/tmrts/flamingo/pkg/flog/logfield"
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
		logfield.Event("main.InitializeContextualization"),
	)
	return nil
}

func FinalizeContextualization() error {
	// TODO(tmrts): Add plug-in hooks
	flog.Info("Finalizing contextualization",
		logfield.Event("main.FinalizeContextualization"),
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
	flag.Parse()

	/*
	 *    if flags.cloudConfig != "" {
	 *        conf, err := os.Open(flags.cloudConfig)
	 *        if err != nil {
	 *            log.Panicf("fatal error config file: %v", err)
	 *        }
	 *
	 *        _, err = cloudconfig.Parse(conf)
	 *        if err != nil {
	 *            log.Panicf("fatal error config file: %v", err)
	 *        }
	 *    }
	 */

	// TODO(tmrts): Build a logger hierarchy
	hasRoot, err := HasRootPrivileges()
	if err != nil {
		flog.Fatal("Failed checking user privileges",
			logfield.Event("main.HasRootPrivileges"),
			logfield.Error(err),
		)
	}

	if !hasRoot {
		flog.Fatal("current user doesn't have root privileges")
	}

	centOS := distro.CentOS(sys.DefaultExecutor)

	if err := InitializeContextualization(); err != nil {
		flog.Fatal("Failed to start contextualization",
			logfield.Event("main.InitializeContextualization"),
			logfield.Error(err),
		)
	}

	providers := datasrc.SupportedProviders()

	p, err := datasrc.FindProvider(providers, 5*time.Second)
	if err != nil {
		flog.Fatal("Failed to start contextualization",
			logfield.Event("datasrc.FindProvider"),
			logfield.Error(err),
		)
	}

	m, err := p.FetchMetadata()
	if err != nil {
		flog.Fatal("Failed to fetch meta-data from provider",
			logfield.Event("%s.FetchMetadata", p),
			logfield.Error(err),
		)
	}

	u, err := p.FetchUserdata()
	if err != nil {
		flog.Fatal("Failed to fetch user-data from provider",
			logfield.Event("%s.FetchUserdata", p),
			logfield.Error(err),
		)
	}

	if err := centOS.ConsumeMetadata(m); err != nil {
		flog.Fatal("Failed to consume meta-data",
			logfield.Event("CentOS.ConsumeMetadata"),
			logfield.Error(err),
		)
	}

	if err := centOS.ConsumeUserdata(u); err != nil {
		flog.Fatal("Failed to consume user-data",
			logfield.Event("CentOS.ConsumeMetadata"),
			logfield.Error(err),
		)
	}

	if err := FinalizeContextualization(); err != nil {
		flog.Fatal("Failed to finalize contextualization",
			logfield.Event("main.FinalizeContextualization"),
			logfield.Error(err),
		)
	}
}
