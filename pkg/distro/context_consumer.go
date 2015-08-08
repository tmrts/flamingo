package distro

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata/cloudconfig"
	"github.com/tmrts/flamingo/pkg/file"
	"github.com/tmrts/flamingo/pkg/flog"
	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/sys/nss"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

type ContextConsumer interface {
	ConsumeUserdata(userdata.Map) error
	ConsumeMetadata(*metadata.Digest) error
}

// ConsumeScript writes the given contents to a temporary file
// and executes the file.
func (imp *Implementation) ConsumeScript(c string) error {
	tempFile := &context.TempFile{
		Content:     c,
		Permissions: 0600,
	}

	return <-context.Using(tempFile, func(f *os.File) error {
		f.Close()
		_, err := imp.Execute("sh", f.Name())

		return err
	})
}

func sTorc(s string) (rc io.ReadCloser) {
	return ioutil.NopCloser(strings.NewReader(s))
}

func (imp *Implementation) consumeCloudConfig(contents string) error {
	conf, err := cloudconfig.Parse(sTorc(contents))
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
			Event: "distro.consumeCloudConfig",
		},
	)

	for _, f := range conf.Files {
		p, err := strconv.Atoi(f.Permissions)
		if err != nil {
			flog.Error("Failed to convert permissions",
				flog.Fields{
					Event: "strconv.Atoi",
					Error: err,
				},
				flog.Details{
					"file":        f.Path,
					"permissions": f.Permissions,
				},
			)
			continue
		}

		perms := os.FileMode(p)

		err = file.New(f.Path, file.Permissions(perms), file.Contents(f.Content))
		if err != nil {
			flog.Error("Failed to create file",
				flog.Fields{
					Event: "file.New",
					Error: err,
				},
				flog.Details{
					"file": f.Path,
				},
			)
		}
	}

	for _, cmd := range conf.Commands {
		flog.Debug("Executing command",
			flog.Fields{
				Event: "distro.consumeCloudConfig",
			},
			flog.Details{
				"command": cmd,
			},
		)

		out, err := imp.Execute(cmd[0], cmd[1:]...)
		if err != nil {
			flog.Error("Failed to execute command",
				flog.Fields{
					Event: "Implementation.Execute",
					Error: err,
				},
				flog.Details{
					"command": cmd,
				},
			)
		}

		flog.Debug("Executed command",
			flog.Fields{
				Event: "identityManager.CreateGroup",
			},
			flog.Details{
				"command": cmd,
				"output":  out,
			},
		)
	}

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

	for userName, sshKeys := range conf.AuthorizedKeys {
		flog.Info("Authorizing SSH keys",
			flog.Fields{
				Event: "distro.consumeCloudConfig",
			},
			flog.Details{
				"user": userName,
			},
		)

		usr, err := nss.GetUser(userName)
		if err != nil {
			flog.Error("Failed to retrieve user NSS entry",
				flog.Fields{
					Event: "nss.GetUser",
					Error: err,
				},
				flog.Details{
					"user": userName,
				},
			)
			continue
		}

		if err := ssh.AuthorizeKeysFor(usr, sshKeys); err != nil {
			flog.Error("Failed to authorize SSH keys for user",
				flog.Fields{
					Event: "ssh.AuthorizeKeysFor",
					Error: err,
				},
				flog.Details{
					"user": userName,
				},
			)
		}
	}

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

	confs := u.CloudConfigs()
	if len(confs) > 1 {
		flog.Info("Found multiple cloud-config files",
			flog.Fields{
				Event: "distro.ConsumeUserdata",
			},
		)
	}

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

	for name, content := range scripts {
		flog.Info("Executing script",
			flog.Fields{
				Event: "distro.ConsumeUserdata",
			},
			flog.Details{
				"name": name,
			},
		)

		if err := imp.ConsumeScript(content); err != nil {
			flog.Error("Failed to execute script",
				flog.Fields{
					Event: "distro.ConsumeScript",
				},
				flog.Details{
					"name": name,
				},
				flog.DebugFields{
					"content": content,
				},
			)
		}
	}

	flog.Info("Finished consuming user-data",
		flog.Fields{
			Event: "distro.ConsumeUserdata",
		},
	)

	return nil
}

// ConsumeUserdata uses the given userdata to contextualize the distribution implementation.
func (imp *Implementation) ConsumeMetadata(m *metadata.Digest) error {
	flog.Info("Consuming meta-data",
		flog.Fields{
			Event: "distro.ConsumeMetadata",
		},
	)

	if err := imp.SetHostname(m.Hostname); err != nil {
		flog.Error("Failed to set hostname",
			flog.Fields{
				Event: "distro.SetHostname",
			},
			flog.Details{
				"name": m.Hostname,
			},
		)
		return err
	}

	for userName, sshKeys := range m.SSHKeys {
		flog.Info("Authorizing SSH keys",
			flog.Fields{
				Event: "distro.consumeMetadata",
			},
			flog.Details{
				"user": userName,
			},
		)

		usr, err := nss.GetUser(userName)
		if err != nil {
			flog.Error("Failed to retrieve user NSS entry",
				flog.Fields{
					Event: "nss.GetUser",
					Error: err,
				},
				flog.Details{
					"user": userName,
				},
			)
			continue
		}

		if err := ssh.AuthorizeKeysFor(usr, sshKeys); err != nil {
			flog.Error("Failed to authorize SSH keys for user",
				flog.Fields{
					Event: "ssh.AuthorizeKeysFor",
					Error: err,
				},
				flog.Details{
					"user": userName,
				},
				flog.DebugFields{
					"SSHKeys": sshKeys,
				},
			)
		}
	}

	flog.Info("Finished consuming meta-data",
		flog.Fields{
			Event: "distro.ConsumeUserdata",
		},
	)

	return nil
}
