package provider

import (
	"fmt"

	"github.com/tmrts/flamingo/pkg/datasrc/metadata"
	"github.com/tmrts/flamingo/pkg/datasrc/userdata"
)

// FormatURL is the meta-data service url as a format string
// to be replaced by Version.
// Example:
//	var u FormatURL = "http://169.256.169.256/metadata/%v/%v"
//	fmt.Printf(u.Fill("latest", "hostname")) -> http://169.256.169.256/metadata/latest/hostname"
type FormatURL string

func (u FormatURL) Fill(values ...interface{}) string {
	return fmt.Sprintf(string(u), values...)
}

// Interface is the interface that wraps the meta-data retrieval method.
type Interface interface {
	metadata.Provider
	userdata.Provider
}
