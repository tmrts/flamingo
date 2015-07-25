package userdata

import (
	"errors"
	"strings"

	"github.com/tmrts/flamingo/pkg/datasrc/userdata/cloudconfig"
)

// Errors regarding the user-data file types.
var (
	ErrNotACloudConfigFile = errors.New("userdata: not a cloud-config file")
	ErrNotAScriptFile      = errors.New("userdata: not a user-data script file")
)

// Map contains the user-data attributes given by the provider.
type Map map[string]string

// IsScript checks whether the given content belongs to a script file.
func IsScript(content string) error {
	if !strings.HasPrefix(content, "#! ") {
		return ErrNotAScriptFile
	}

	return nil
}

// Scripts return a new map containing only the contents
// that are valid scripts.
func (m Map) Scripts() map[string]string {
	scripts := make(map[string]string)

	for k, v := range m {
		if err := IsScript(v); err == nil {
			scripts[k] = v
		}
	}

	return scripts
}

// IsCloudConfig checks whether the given content belongs to a cloud-config file.
func IsCloudConfig(content string) error {
	if !strings.HasPrefix(content, "#cloud-config\n") {
		return cloudconfig.ErrNotACloudConfigFile
	}

	return nil
}

// CloudConfigs return a new map containing only the contents
// that are valid cloud-config files.
func (m Map) CloudConfigs() map[string]string {
	confs := make(map[string]string)

	for k, v := range m {
		if err := IsCloudConfig(v); err == nil {
			confs[k] = v
		}
	}

	return confs
}

// Provider retrieves user-data attributes from the meta-data server
// and returns a Map.
type Provider interface {
	FetchUserdata() (Map, error)
}
