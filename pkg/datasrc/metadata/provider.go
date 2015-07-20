package metadata

import "fmt"

// FormatURL is the meta-data service url as a format string
// to be replaced by Version.
// Example:
//	var u FormatURL = "http://169.256.169.256/metadata/%v/%v"
//	fmt.Printf(u.Fill("latest", "hostname")) -> http://169.256.169.256/metadata/latest/hostname"
type FormatURL string

func (u FormatURL) Fill(values ...interface{}) string {
	return fmt.Sprintf(string(u), values...)
}

// Provider is the interface that wraps the meta-data retrieval method.
type Provider interface {
	MetaData(Version) (Interface, error)
}
