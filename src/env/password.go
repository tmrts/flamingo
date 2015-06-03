package env

type Password struct {
	Hash string
}

func (p Password) String() string {
	return p.Hash
}
