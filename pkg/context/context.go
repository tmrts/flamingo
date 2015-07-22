package context

import (
	"os"

	"github.com/tmrts/flamingo/pkg/file"
)

type Interface interface {
	Enter() error
	Exit() error

	Use(interface{}) error
}

type File struct {
	source *os.File

	Path        string
	Permissions os.FileMode
}

func (f *File) Enter() error {
	src, err := os.OpenFile(f.Path, os.O_WRONLY, f.Permissions)
	f.source = src
	return err
}

func (f *File) Exit() error {
	return f.source.Close()
}

// Expects func(*os.File) error
func (f *File) Use(fn interface{}) error {
	closure := fn.(func(*os.File) error)

	return closure(f.source)
}

type NewFile struct {
	source *os.File

	Path        string
	Permissions os.FileMode
}

func (nf *NewFile) Enter() error {
	src, err := os.OpenFile(nf.Path, os.O_CREATE|os.O_WRONLY, nf.Permissions)
	nf.source = src
	return err
}

func (nf *NewFile) Exit() error {
	return nf.source.Close()
}

// Expects func(*os.File) error
func (nf *NewFile) Use(fn interface{}) error {
	closure := fn.(func(*os.File) error)

	return closure(nf.source)
}

type TempFile struct {
	source *os.File
	fname  string

	Content     string
	Permissions os.FileMode
}

func (tf *TempFile) Enter() error {
	if tf.Permissions == 0 {
		tf.Permissions = os.FileMode(0600)
	}

	src, err := file.Temp(tf.Content, tf.Permissions)
	tf.source = src
	tf.fname = src.Name()
	return err
}

func (tf *TempFile) Exit() error {
	tf.source.Close()
	return os.Remove(tf.fname)
}

// Expects func(*os.File) error
func (tf *TempFile) Use(fn interface{}) error {
	closure := fn.(func(*os.File) error)

	return closure(tf.source)
}

func Using(ctxt Interface, fn interface{}) chan error {
	errch := make(chan error, 1)

	go func(chan error) {
		err := ctxt.Enter()
		if err != nil {
			errch <- err
			return
		} else {
			defer ctxt.Exit()
		}

		errch <- ctxt.Use(fn)
	}(errch)

	return errch
}
