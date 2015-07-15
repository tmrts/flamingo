package metadata

import "fmt"

type Version string

type URLFormat string

func (u URLFormat) WithVersion(v Version) string {
	return fmt.Sprintf(string(u), v)
}

type Provider interface {
	MetaData(Version) (Interface, error)
}
