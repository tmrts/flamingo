// Package strutil provides different string utilities
package strutil

import (
	"io"
	"io/ioutil"
	"strings"
)

func ToReader(s string) (r io.Reader) {
	return strings.NewReader(s)
}

func ToReadCloser(s string) (rc io.ReadCloser) {
	return ioutil.NopCloser(ToReader(s))
}
