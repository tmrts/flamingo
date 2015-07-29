package main

import (
	"flag"
	"log"
	"time"

	"github.com/tmrts/flamingo/pkg/datasrc"
	"github.com/tmrts/flamingo/pkg/distro"
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
	log.Print("main: starting the contextualization")
	return nil
}

func FinalizeContextualization() error {
	// TODO(tmrts): Add plug-in hooks
	log.Print("main: finalizing the contextualization")
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
		log.Fatalf("failed to check current user privileges, %v", err)
	}

	if !hasRoot {
		log.Fatal("current user doesn't have root privileges")
	}

	centOS := distro.CentOS(sys.DefaultExecutor)

	if err := InitializeContextualization(); err != nil {
		log.Fatalf("failed to contextualization, %v", err)
	}

	providers := datasrc.SupportedProviders()

	p, err := datasrc.FindProvider(providers, 5*time.Second)
	if err != nil {
		log.Fatalf("failed to find an available datasource provider, %v", err)
	}

	m, err := p.FetchMetadata()
	if err != nil {
		log.Fatalf("failed to fetch meta-data from provider %v, %v", p, err)
	}

	u, err := p.FetchUserdata()
	if err != nil {
		log.Fatalf("failed to fetch user-data from provider %v, %v", p, err)
	}

	if err := centOS.ConsumeMetadata(m); err != nil {
		log.Fatalf("failed to consume meta-data, %v", err)
	}

	if err := centOS.ConsumeUserdata(u); err != nil {
		log.Fatalf("failed to consume user-data, %v", err)
	}

	if err := FinalizeContextualization(); err != nil {
		log.Fatalf("failed to finalize contextualization, %v", err)
	}
}
