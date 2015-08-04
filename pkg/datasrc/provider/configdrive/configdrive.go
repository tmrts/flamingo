// Handling configuration data passed via Config-Drive disk images.
package configdrive

import (
	"errors"

	"github.com/tmrts/flamingo/pkg/sys"
)

var (
	ErrUnableToLocateConfigDrive = errors.New("configdrive: couldn't locate the config drive device")
)

// FindMountTarget finds the mount target location of the config drive device.
// It finds the device labeled with 'LABEL=config-2', then gets the mount target
// of that device.
func FindMountTarget(e sys.Executor) (string, error) {
	dev, err := e.Execute("blkid", "-t LABEL='config-2'", "-odevice")
	if err != nil {
		return "", err
	}

	if dev == "" {
		return "", ErrUnableToLocateConfigDrive
	}

	return e.Execute("findmnt", "--raw", "--noheadings", "--output TARGET", dev)
}

/*
 *
 *func (m *Mount) FetchMetadata() (metadata.Digest, error) {
 *    var data []byte
 *    var m struct {
 *        SSHKeys map[string][]string `json:"public_keys"`
 *        Hostname            string            `json:"hostname"`
 *        NetworkInterfaces   metadata.NetworkInterfaces `json:"content_path"`
 *    }
 *
 *    metadata.SSHPublicKeys = m.SSHKeys
 *    metadata.Hostname = m.Hostname
 *    if m.NetworkConfig.ContentPath != "" {
 *        metadata.NetworkConfig, err = m.tryReadFile(path.Join(m.Path, m.NetworkConfig.ContentPath))
 *    }
 *
 *    return metadata, err
 *}
 */
