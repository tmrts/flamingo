// Package cloudconfig validates and parses a cloud-config data file.
package cloudconfig

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/sys/ssh"

	"gopkg.in/yaml.v2"
)

// WriteFile represents the write_file directive found in a
// cloud-config file. It contains the fields necessary to create
// a new file in the system.
type WriteFile struct {
	Path        string
	Owner       string
	Permissions string
	Encoding    string

	Content string
}

var (
	ErrNotACloudConfigFile = errors.New("cloudconfig: not a cloud-config file")
)

// IsValid reads from the given io.ReadCloser and determines
// whether the read contents belong to a valid cloud-config file.
// It returns an ErrNotACloudConfigFile if the file is not a valid cloud-config file.
// If it encounters an error while reading from the io.ReadCloser it returns the error.
func IsValid(rdr io.ReadCloser) error {
	buf, err := ioutil.ReadAll(rdr)
	if err != nil {
		return err
	}
	rdr.Close()

	if contents := string(buf); strings.HasPrefix(contents, "#cloud-config\n") != true {
		return ErrNotACloudConfigFile
	}

	return nil
}

type user struct {
	identity.User     `yaml:""`
	AuthorizedSSHKeys []string `yaml:"ssh-authorized-keys"`
	SSHImportID       string   `yaml:"ssh-import-id"`
}

type cloudConfig struct {
	RunCMD         []interface{}
	AuthorizedKeys []string          `yaml:"ssh_authorized_keys"`
	SSHKeyPairs    map[string]string `yaml:"ssh_keys"`
	Users          []identity.User   `yaml:"users"`
	Groups         []interface{}     `yaml:"groups"`
	Files          []WriteFile       `yaml:"write_files"`
}

// Digest is the parsed cloud-config file.
type Digest struct {
	Commands       [][]string
	Files          []WriteFile
	Groups         map[string][]string
	Users          map[string]identity.User
	AuthorizedKeys map[string][]ssh.Key
	SSHKeyPairs    []ssh.KeyPair
}

// Parse reads from the given io.ReadCloser and
// parses read contents of a cloud-config file.
func Parse(rdr io.ReadCloser) (*Digest, error) {
	buf, err := ioutil.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	rdr.Close()

	var conf cloudConfig
	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}
	var c Digest

	c.Commands = parseCommands(conf.RunCMD)

	c.AuthorizedKeys = make(map[string][]ssh.Key)
	for _, key := range conf.AuthorizedKeys {
		c.AuthorizedKeys["root"] = append(c.AuthorizedKeys["root"], ssh.Key(key))
	}

	public_keys, private_keys := make(map[string]string), make(map[string]string)

	// TODO: Extend ssh key syntax beyond cloud-config
	for k, v := range conf.SSHKeyPairs {
		if strings.HasSuffix(k, "private") {
			private_keys[strings.TrimSuffix(k, "_private")] = v
		} else {
			public_keys[strings.TrimSuffix(k, "_public")] = v
		}
	}

	for key, value := range public_keys {
		c.SSHKeyPairs = append(c.SSHKeyPairs, ssh.KeyPair{
			Public:  ssh.Key(value),
			Private: ssh.Key(private_keys[key]),
		})
	}

	c.Groups = parseGroups(conf.Groups)
	c.Users = parseUsers(conf.Users)
	/*
	* BUG(yaml.v2): Embedded structs are not unmarshaled properly.
	* TODO(tmrts): Use another yaml library or extend identity.User.
	*for _, usr := range conf.Users {
	*    for _, key := range usr.AuthorizedSSHKeys {
	*       c.AuthorizedKeys[usr.Name] = append(c.AuthorizedKeys[usr.Name], ssh.Key(key))
	*    }
	*}
	 */

	c.Files = conf.Files

	return &c, nil
}

func parseCommands(runcmd []interface{}) (commands [][]string) {
	var cmds []string

	for _, cmd := range runcmd {
		switch cmd := cmd.(type) {
		case string:
			cmds = []string{"sh", "-c", cmd}
		case []interface{}:
			cmds = []string{}
			for _, s := range cmd {
				cmds = append(cmds, s.(string))
			}
		}

		commands = append(commands, cmds)
	}

	return
}

func parseUsers(users []identity.User) map[string]identity.User {
	userMap := make(map[string]identity.User)

	for _, user := range users {
		userMap[user.Name] = user
	}

	return userMap
}

func parseGroups(vals []interface{}) map[string][]string {
	groups := make(map[string][]string)

	for _, group := range vals {
		switch group := group.(type) {
		case string:
			groups[group] = []string{}
		case map[interface{}]interface{}:
			for name, users := range group {
				name := name.(string)

				members := []string{}
				for _, m := range users.([]interface{}) {
					members = append(members, m.(string))
				}

				groups[name] = members
			}
		}
	}

	return groups
}
