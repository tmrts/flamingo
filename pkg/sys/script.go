package sys

import "io/ioutil"

// HasShabang checks whether the named file is a script by checking the shabang directive
func FileHasShabang(filename string) (bool, error) {
	// TODO: Don't read the complete file, only the first line or few bytes.
	byteContent, _ := ioutil.ReadFile(filename)

	content := string(byteContent)

	return content[0] == '#' && content[1] == '!', nil
}

// WriteScript writes the content of the byte slice into the named file
func WriteScript(filename string, contents []byte) error {
	// TODO: Variable environment preparation (e.g. temporarily Chdir)
	// BUG: The contents are appended if file/script already exists.
	//		Arbitrary code could be execution can be performed.
	return ioutil.WriteFile(filename, contents, 0744)
}
