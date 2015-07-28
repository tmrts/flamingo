package distro

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tmrts/flamingo/pkg/context"
	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata/cloudconfig"
	"github.com/tmrts/flamingo/pkg/file"
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
		out, err := imp.Execute("sh", f.Name())

		log.Printf("ConsumeScript script output -> %v", out)

		return err
	})
}

func sTorc(s string) (rc io.ReadCloser) {
	return ioutil.NopCloser(strings.NewReader(s))
}

func (imp *Implementation) consumeCloudConfig(contents string) error {
	conf, err := cloudconfig.Parse(sTorc(contents))
	if err != nil {
		log.Printf("consumeCloudConfig -> consuming %v,  %v", err)
		return err
	}

	for _, f := range conf.Files {
		log.Printf("consumeCloudConfig -> writing files %v", f.Path)

		p, err := strconv.Atoi(f.Permissions)
		if err != nil {
			log.Printf("consumeCloudConfig -> error writing files %v", f.Path, err)
			continue
		}

		perms := os.FileMode(p)

		err = file.New(f.Path, file.Permissions(perms), file.Contents(f.Content))
		if err != nil {
			log.Printf("consumeCloudConfig -> error writing files %v", f.Path, err)
		}
	}

	for _, cmd := range conf.Commands {
		log.Printf("consumeCloudConfig -> executing %v", cmd)

		out, err := imp.Execute(cmd[0], cmd[1:]...)
		if err != nil {
			log.Printf("consumeCloudConfig -> error executing %v, %v", cmd, err)
		}

		log.Printf("consumeCloudConfig command output -> %v", out)

	}

	for grpName, _ := range conf.Groups {
		log.Printf("consumeCloudConfig -> creating group %v", grpName)

		newGrp := identity.Group{
			Name: grpName,
		}

		if err := imp.ID.CreateGroup(newGrp); err != nil {
			log.Printf("consumeCloudConfig -> error creating group %v, %v", grpName, err)
		}
	}

	for _, usr := range conf.Users {
		log.Printf("consumeCloudConfig -> creating group %v", usr.Name)

		if err := imp.ID.CreateUser(usr); err != nil {
			log.Printf("consumeCloudConfig -> error creating user %v, %v", usr.Name, err)
		}
	}

	for grpName, usrNames := range conf.Groups {
		for _, usrName := range usrNames {
			log.Printf("consumeCloudConfig -> adding user %v to group %v", usrName, grpName)

			if err := imp.ID.AddUserToGroup(usrName, grpName); err != nil {
				log.Printf("consumeCloudConfig -> error adding user %v to group %v, %v", usrName, grpName, err)
			}
		}
	}

	for userName, sshKeys := range conf.AuthorizedKeys {
		log.Printf("consumeCloudConfig -> authorizing ssh keys for user %v", userName)

		usr, err := nss.GetUser(userName)
		if err != nil {
			log.Printf("consumeCloudConfig -> error retrieving entry for user %v, %v", userName, err)
			continue
		}

		log.Printf("consumeCloudConfig -> authorizing ssh keys for user %v", userName)
		if err := ssh.AuthorizeKeysFor(usr, sshKeys); err != nil {
			log.Printf("consumeCloudConfig -> error authorizing ssh keys for user %v, %v", userName, err)
		}
	}

	return err
}

// ConsumeUserdata uses the given userdata to contextualize the distribution implementation.
func (imp *Implementation) ConsumeUserdata(u userdata.Map) error {
	// TODO(tmrts): Store unused user-data in files?
	// TODO(tmrts): Execute scripts in rc.local or a similar level
	log.Print("ConsumeUserdata -> entering")

	// TODO(tmrts): Use only scripts with 'startup', 'shutdown', 'user-data'.
	scripts := u.Scripts()

	confs := u.CloudConfigs()
	if len(confs) > 1 {
		log.Print("ConsumeUserdata -> more than one cloud-config file found")
	}

	for name, content := range confs {
		log.Printf("ConsumeUserdata -> consuming %v", name)

		err := imp.consumeCloudConfig(content)
		if err != nil {
			return err
		}
	}

	for name, content := range scripts {
		log.Printf("ConsumeUserdata -> executing script %v", name)

		if err := imp.ConsumeScript(content); err != nil {
			log.Printf("ConsumeUserdata -> error executing script %v, %v", name, err)
		}
	}

	log.Print("ConsumeUserdata -> exiting")

	return nil
}

// ConsumeUserdata uses the given userdata to contextualize the distribution implementation.
func (imp *Implementation) ConsumeMetadata(m *metadata.Digest) error {
	log.Print("ConsumeMetadata -> entering")

	if err := imp.SetHostname(m.Hostname); err != nil {
		return err
	}

	for userName, sshKeys := range m.SSHKeys {
		log.Printf("ConsumeMetadata -> authorizing ssh keys for user %v", userName)

		usr, err := nss.GetUser(userName)
		if err != nil {
			log.Printf("ConsumeMetadata -> error retrieving user from nss, %v", err)
			continue
		}

		if err := ssh.AuthorizeKeysFor(usr, sshKeys); err != nil {
			log.Printf("consumeCloudConfig -> error authorizing ssh keys for user %v, %v", userName, err)
		}
	}

	log.Print("ConsumeMetadata -> exiting")

	return nil
}
