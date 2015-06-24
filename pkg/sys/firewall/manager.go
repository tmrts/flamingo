package firewall

type Table string

const (
	Filter Table = "filter"
	Nat    Table = "nat"
)

type Chain struct {
	Name  string
	Table Table
}

type Manager interface {
	CheckDependencies() error
	AppendTo(Chain, string) error
	DeleteFrom(Chain, string) error
	InsertTo(Chain, string) error
	Check(Chain, string) (bool, error)
}
