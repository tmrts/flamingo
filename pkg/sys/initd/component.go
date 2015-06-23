package initd

type Component interface {
	Name() string
	Path() string
}
