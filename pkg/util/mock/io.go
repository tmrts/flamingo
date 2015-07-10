package mock

import "bytes"

type ReadCloser struct {
	*bytes.Buffer
	IsClosed bool
}

func (mrc *ReadCloser) Close() (err error) {
	mrc.IsClosed = true
	return
}

func NewReadCloser(contents string) *ReadCloser {
	return &ReadCloser{bytes.NewBufferString(contents), false}
}
