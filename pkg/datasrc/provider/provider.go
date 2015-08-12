// Package provider contains the interface for a data source provider
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

// Fill method is a wrapper around fmt.Sprintf to fill format URLs.
func (u FormatURL) Fill(values ...interface{}) string {
	return fmt.Sprintf(string(u), values...)
}

// Interface is the interface that represents the ability to
// provide both user-data and meta-data information.
type Interface interface {
	metadata.Provider
	userdata.Provider
}
