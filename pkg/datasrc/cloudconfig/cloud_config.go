package cloudconfig

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/tmrts/flamingo/pkg/sys/identity"
	"github.com/tmrts/flamingo/pkg/sys/ssh"

	"gopkg.in/yaml.v2"
)

type cloudConfig struct {
	RunCMD         []interface{}
	AuthorizedKeys []string            `yaml:"ssh_authorized_keys"`
	SSHKeyPairs    map[string]string   `yaml:"ssh_keys"`
	Users          []interface{}       `yaml:"users"`
	Groups         map[string][]string `yaml:"groups"`
}

type Contextualization struct {
	Commands       []string
	AuthorizedKeys []ssh.Key
	SSHKeyPairs    []ssh.KeyPair
	Groups         map[string][]string
	Users          []identity.User
}

func Parse(filename string) (*Contextualization, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var conf cloudConfig
	if err := yaml.Unmarshal(f, &conf); err != nil {
		return nil, err
	}
	var c Contextualization
	for _, cmd := range conf.RunCMD {
		switch cmd.(type) {
		case string:
			c.Commands = append(c.Commands, cmd.(string))
		case []interface{}:
			var commands []string
			for _, s := range cmd.([]interface{}) {
				commands = append(commands, s.(string))
			}

			c.Commands = append(c.Commands, strings.Join(commands, " "))
		}
	}

	for _, key := range conf.AuthorizedKeys {
		c.AuthorizedKeys = append(c.AuthorizedKeys, []byte(key))
	}

	public_keys, private_keys := make(map[string]string), make(map[string]string)

	// TODO: Extend ssh key syntax beyond cloud-config
	for k, v := range conf.SSHKeyPairs {
		if strings.HasSuffix("private", k) {
			private_keys[strings.TrimSuffix(k, "_private")] = v
		} else {
			public_keys[strings.TrimSuffix(k, "_public")] = v
		}
	}

	for key, value := range public_keys {
		c.SSHKeyPairs = append(c.SSHKeyPairs, ssh.KeyPair{
			Public:  []byte(value),
			Private: []byte(private_keys[key]),
		})
	}

	for group, users := range conf.Groups {
		fmt.Printf("%v: %v\n", group, users)
	}

	return &c, nil
}
