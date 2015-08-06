// Handling configuration data passed via Config-Drive disk images.
package configdrive

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/provider/openstack"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
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

type Mount struct {
	Path string
}

func (m Mount) FetchMetadata() (*metadata.Digest, error) {
	metadataPath := filepath.Join(m.Path, "openstack/2012-08-10/meta_data.json")
	buf, err := ioutil.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var metadata openstack.Metadata
	if err := json.Unmarshal(buf, &metadata); err != nil {
		return nil, err
	}

	digest := metadata.Digest()

	return &digest, err
}

func (m Mount) FetchUserdata() (userdata.Map, error) {
	userdataPath := filepath.Join(m.Path, "openstack/2012-08-10/user_data")
	buf, err := ioutil.ReadFile(userdataPath)
	if err != nil {
		return nil, err
	}

	userdata := make(userdata.Map)

	userdata["user-data"] = string(buf)

	return userdata, err
}
