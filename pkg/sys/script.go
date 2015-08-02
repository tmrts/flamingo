package sys

import (
	"io/ioutil"
	"os"
)

// HasShabang checks whether the named file is a script by checking the shabang directive
func FileHasShabang(filename string) (bool, error) {
	f, err := os.Open(filename)
	if err != nil {
		return false, err
	} else {
		defer f.Close()
	}

	twoBytes := make([]byte, 2)

	_, err = f.Read(twoBytes)
	if err != nil {
		return false, err
	}

	return twoBytes[0] == '#' && twoBytes[1] == '!', nil
}

// WriteScript writes the content of the byte slice into the named file
func WriteScript(filename string, contents []byte) error {
	// TODO: Variable environment preparation (e.g. temporarily Chdir)
	// BUG: The contents are appended if file/script already exists.
	//		Arbitrary code could be execution can be performed.
	return ioutil.WriteFile(filename, contents, 0744)
}
