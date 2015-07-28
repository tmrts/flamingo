package file

type System interface {
	New(string, ...argument) error
	WriteTo(string, ...argument) error
}
