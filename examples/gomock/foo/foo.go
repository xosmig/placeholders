package foo

type Foo struct {
	Bar *Bar
	Baz BazPtrWrapper
}

func NewFoo(bar *Bar, baz BazPtrWrapper) *Foo {
	return &Foo{bar, baz}
}

type Bar struct {
	S string
}

func NewBar(s string) *Bar {
	return &Bar{s}
}

type Baz struct {
	S string
}

func NewBaz(s string) *Baz {
	return &Baz{s}
}

type BazPtrWrapper *Baz

//go:generate mockgen -destination ./mock/foo_processor.mock.go . FooProcessor
type FooProcessor interface {
	Process(foo *Foo) string
}
