package configdrive

import (
	"encoding/json"
	"os"
	"path"
)

type Mount struct {
	Path string
}

func (m *Mount) Available() bool {
	_, err := os.Stat(m.path)
	return !os.IsNotExist(err)
}

func (m *Mount) Metadata() (metadata.Digest, error) {
	var data []byte
	var m struct {
		SSHAuthorizedKeyMap map[string]string `json:"public_keys"`
		Hostname            string            `json:"hostname"`
		NetworkConfig       struct {
			ContentPath string `json:"content_path"`
		} `json:"network_config"`
	}

	if data, err = m.tryReadFile(path.Join(m.Path, "meta_data.json")); err != nil || len(data) == 0 {
		return
	}
	if err = json.Unmarshal([]byte(data), &m); err != nil {
		return
	}

	metadata.SSHPublicKeys = m.SSHAuthorizedKeyMap
	metadata.Hostname = m.Hostname
	if m.NetworkConfig.ContentPath != "" {
		metadata.NetworkConfig, err = m.tryReadFile(path.Join(m.Path, m.NetworkConfig.ContentPath))
	}

	userF := m.tryReadFile(path.Join(m.Path, "user_data"))

	return metadata.Combine(userF)
}
