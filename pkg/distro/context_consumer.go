package distro

import (
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata/cloudconfig"
	"github.com/tmrts/flamingo/pkg/flog"
	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/util/strutil"
)

// ContextConsumer is the interface that represents the ability to
// use a metadata.Digest and a userdata.Map for mutating itself.
// The expected implementers of the interface are operating
// system implementations.
type ContextConsumer interface {
	ConsumeUserdata(userdata.Map) error
	ConsumeMetadata(*metadata.Digest) error
}

// consumeCloudConfig parses the given cloud config file contents and
// consumes the parsed directives.
func (imp *Implementation) consumeCloudConfig(contents string) error {
	conf, err := cloudconfig.Parse(strutil.ToReadCloser(contents))
	if err != nil {
		flog.Error("Failed to parse cloud config file",
			flog.Fields{
				Event: "cloudconfig.Parse",
				Error: err,
			},
		)

		return err
	}

	flog.Debug("Persisting files",
		flog.Fields{
			Event: "distro.Implementation.consumeCloudConfig",
		},
	)

	imp.writeFiles(conf.Files)

	imp.consumeCommands(conf.Commands)

	for grpName, _ := range conf.Groups {
		flog.Info("Creating user group",
			flog.Fields{
				Event: "distro.consumeCloudConfig",
			},
			flog.Details{
				"group": grpName,
			},
		)

		newGrp := identity.Group{
			Name: grpName,
		}

		if err := imp.ID.CreateGroup(newGrp); err != nil {
			flog.Error("Failed to create a user group",
				flog.Fields{
					Event: "identityManager.CreateGroup",
					Error: err,
				},
				flog.Details{
					"group": grpName,
				},
			)
		}
	}

	for _, usr := range conf.Users {
		flog.Info("Creating user",
			flog.Fields{
				Event: "distro.consumeCloudConfig",
			},
			flog.Details{
				"user": usr.Name,
			},
		)

		if err := imp.ID.CreateUser(usr); err != nil {
			flog.Error("Failed to create a user",
				flog.Fields{
					Event: "identityManager.CreateUser",
					Error: err,
				},
				flog.Details{
					"user": usr.Name,
				},
			)
		}
	}

	for grpName, usrNames := range conf.Groups {
		for _, usrName := range usrNames {
			flog.Info("Adding user to group",
				flog.Fields{
					Event: "distro.consumeCloudConfig",
				},
				flog.Details{
					"user":  usrName,
					"group": grpName,
				},
			)

			if err := imp.ID.AddUserToGroup(usrName, grpName); err != nil {
				flog.Error("Failed to add user to group",
					flog.Fields{
						Event: "identityManager.AddUserToGroup",
						Error: err,
					},
					flog.Details{
						"user":  usrName,
						"group": grpName,
					},
				)
			}
		}
	}

	imp.consumeSSHKeys(conf.AuthorizedKeys)

	return err
}

// ConsumeUserdata uses the given userdata to contextualize the distribution implementation.
func (imp *Implementation) ConsumeUserdata(u userdata.Map) error {
	// TODO(tmrts): Store unused user-data in files?
	// TODO(tmrts): Execute scripts in rc.local or a similar level
	flog.Info("Consuming user-data",
		flog.Fields{
			Event: "distro.ConsumeUserdata",
		},
	)

	// TODO(tmrts): Use only scripts with 'startup', 'shutdown', 'user-data'.
	scripts := u.Scripts()
	flog.Info("Searched for script files",
		flog.Fields{
			Event: "userdata.Map.Scripts",
		},
		flog.Details{
			"count": len(scripts),
		},
	)

	imp.consumeScripts(scripts)

	confs := u.CloudConfigs()
	flog.Info("Searched for cloud-config files",
		flog.Fields{
			Event: "userdata.Map.CloudConfigs",
		},
		flog.Details{
			"count": len(confs),
		},
	)

	for name, content := range confs {
		flog.Info("Consuming user-data file",
			flog.Fields{
				Event: "distro.ConsumeUserdata",
			},
			flog.Details{
				"name": name,
			},
		)

		err := imp.consumeCloudConfig(content)
		if err != nil {
			flog.Error("Failed to consume cloud-config file",
				flog.Fields{
					Event: "distro.consumeCloudConfig",
				},
				flog.Details{
					"name": name,
				},
				flog.DebugFields{
					"content": content,
				},
			)
			return err
		}
	}

	flog.Info("Finished consuming user-data",
		flog.Fields{
			Event: "distro.Implementation.ConsumeUserdata",
		},
	)

	return nil
}

// ConsumeMetadata uses the given userdata to contextualize the distribution implementation.
func (imp *Implementation) ConsumeMetadata(m *metadata.Digest) error {
	flog.Info("Consuming meta-data",
		flog.Fields{
			Event: "distro.Implementation.ConsumeMetadata",
		},
	)

	if err := imp.setHostname(m.Hostname); err != nil {
		flog.Error("Failed to set hostname",
			flog.Fields{
				Event: "distro.Implementation.setHostname",
			},
			flog.Details{
				"name": m.Hostname,
			},
		)
		return err
	}

	imp.consumeSSHKeys(m.SSHKeys)

	flog.Info("Finished consuming meta-data",
		flog.Fields{
			Event: "distro.Implementation.ConsumeUserdata",
		},
	)

	return nil
}
