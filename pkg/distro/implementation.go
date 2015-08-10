package distro

import (
	"os"
	"strconv"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata/cloudconfig"
	"github.com/tmrts/flamingo/pkg/file"
	"github.com/tmrts/flamingo/pkg/flog"
	"github.com/tmrts/flamingo/pkg/sys"
	"github.com/tmrts/flamingo/pkg/sys/firewall/iptables"
	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/sys/initd"
	"github.com/tmrts/flamingo/pkg/sys/nss"
	"github.com/tmrts/flamingo/pkg/sys/ssh"
)

// Implementation contains the system Managers that are
// used during the instance contextualization process.
type Implementation struct {
	sys.Executor

	ID    identity.Manager
	NSS   nss.Service
	Initd initd.Manager

	// TODO(tmrts): Change to firewall.Manager interface
	Firewall *iptables.Implementation
}

// SetHostname changes the hostname of the given distribution.
func (imp *Implementation) setHostname(hostname string) error {
	_, err := imp.Execute("hostnamectl", "set-hostname", hostname)

	return err
}

// SetTimeZone changes the hostname of the given distribution.
func (imp *Implementation) setTimeZone(timezone string) error {
	_, err := imp.Execute("timedatectl", "set-timezone", timezone)

	return err
}

// executeScript writes the given contents to a temporary file and executes the file.
func (imp *Implementation) executeScript(c string) error {
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

// consumeScripts executes the given name/script content pairs.
func (imp *Implementation) consumeScripts(scripts map[string]string) {
	for name, content := range scripts {
		flog.Info("Executing script",
			flog.Fields{
				Event: "distro.consumeScripts",
			},
			flog.Details{
				"name": name,
			},
		)

		if err := imp.executeScript(content); err != nil {
			flog.Error("Failed to execute script",
				flog.Fields{
					Event: "distro.executeScript",
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
}

// consumeCommands executes the given command and it's arguments.
func (imp *Implementation) consumeCommands(commands [][]string) {
	for _, cmd := range commands {
		flog.Debug("Executing command",
			flog.Fields{
				Event: "distro.Implementation.consumeCommands",
			},
			flog.Details{
				"command": cmd,
			},
		)

		out, err := imp.Execute(cmd[0], cmd[1:]...)
		if err != nil {
			flog.Error("Failed to execute command",
				flog.Fields{
					Event: "distro.Implementation.Execute",
					Error: err,
				},
				flog.Details{
					"command": cmd,
				},
			)
		}

		flog.Debug("Executed command",
			flog.Fields{
				Event: "distro.Implementation.consumeCommands",
			},
			flog.Details{
				"command": cmd,
				"output":  out,
			},
		)
	}
}

// consumeSSHKeys authorizes the given SSH keys for the user.
func (imp *Implementation) consumeSSHKeys(userKeys map[string][]ssh.Key) {
	for userName, sshKeys := range userKeys {
		flog.Info("Authorizing SSH keys",
			flog.Fields{
				Event: "distro.consumeSSHKeys",
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
}

// writeFiles persists the given files to the system.
func (imp *Implementation) writeFiles(files []cloudconfig.WriteFile) {
	for _, f := range files {
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
}
