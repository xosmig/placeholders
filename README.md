[![Go Reference](https://pkg.go.dev/badge/github.com/xosmig/placeholders.svg)](https://pkg.go.dev/github.com/xosmig/placeholders)
[![Go Report Card](https://goreportcard.com/badge/github.com/xosmig/placeholders)](https://goreportcard.com/report/github.com/xosmig/placeholders)
[![Test](https://github.com/xosmig/placeholders/actions/workflows/test.yml/badge.svg)](https://github.com/xosmig/placeholders/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/xosmig/placeholders/branch/main/graph/badge.svg)](https://codecov.io/gh/xosmig/placeholders)

# placeholders
Placeholders for struct fields for efficient testing with [go-cmp](https://github.com/google/go-cmp).  
The package is non-intrusive, i.e., the original structs do not need to be modified.

For more detailed documentation, see: [placeholders.go](https://github.com/xosmig/placeholders/blob/main/placeholders.go)  
For more examples, see: [placeholders_test.go](https://github.com/xosmig/placeholders/blob/main/placeholders_test.go)

## Example of integration with go-mock:

file [examples/gomock/gomock_example_test.go](https://github.com/xosmig/placeholders/blob/main/examples/gomock/gomock_example_test.go)
```go
func TestGomockExample(t *testing.T) {
	ctrl := gomock.NewController(t)

	mock := mock_foo.NewMockFooProcessor(ctrl)

	// Matches a Foo object with any first argument and second argument equal to foo.NewBaz("world").
	mock.EXPECT().Process(matchers.DiffEq(foo.NewFoo(placeholders.Make[*foo.Bar](t), foo.NewBaz("world")))).
		Return("greetings!")

	// Matches a Foo object the first argument equal to foo.NewBar("goodbye") and any second argument.
	mock.EXPECT().Process(matchers.DiffEq(foo.NewFoo(foo.NewBar("goodbye"), placeholders.Make[foo.BazPtrWrapper](t)))).
		Return("farewell!")

	assert.Equal(t, mock.Process(foo.NewFoo(foo.NewBar("hello"), foo.NewBaz("world"))), "greetings!")

	assert.Equal(t, mock.Process(foo.NewFoo(foo.NewBar("goodbye"), foo.NewBaz("world"))), "farewell!")
}
```

file [examples/gomock/foo/foo.go](https://github.com/xosmig/placeholders/blob/main/examples/gomock/foo/foo.go)
```go
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
```
