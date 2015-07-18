package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/tmrts/flamingo/pkg/datasrc/cloudconfig"
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/identity"
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

func StartContextualization() {
	// TODO: Add plug-in hooks
}

func FinalizeContextualization() {
	// TODO: Add plug-in hooks
}

func main() {
	// TODO: Build Meaningful Loggers
	flag.Parse()

	// cloudconfig
	if flags.cloudConfig != "" {
		cloudConfigContext, err := cloudconfig.Parse(flags.cloudConfig)
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %v", err))
		}
	}

	metadataDigest := metadata.Get(10 * time.Second)

	conf := datasrc.Merge(metadataDigest, cloudConfigContext)

	StartContextualization()

	// usergroups
	idm := identity.Manager{Exec: sys.DefaultExecutor}

	for _, grp := range conf.Groups {
		if err := idm.CreateGroup(grp); err != nil {
			panic(err)
		}
	}

	for _, usr := range conf.Users {
		if err := idm.CreateUser(usr); err != nil {
			panic(err)
		}
	}

	for _, grp := range conf.Groups {
		for _, usr := range grp {
			if err := idm.AddUserToGroup(grp, usr); err != nil {
				panic(err)
			}
		}
	}

	// ssh_keys
	if err := ssh.InitializeFor("root"); err != nil {
		panic(err)
	}

	ssh.AuthorizeSSHKey(f, conf.AuthorizedKeys...)

	// write_files
	for _, f := range conf.Files {
		file.New(f.Name, file.Contents(f.Data), file.Uid(0), file.Gid(0), file.Permissions(f.Perms))
	}

	// run_cmd
	for _, cmd := range conf.Commands {
		sys.Execute(cmd)
	}

	if err := FinalizeContextualization(); err != nil {
		panic(err)
	}
}
