package placeholders

import (
	"github.com/google/go-cmp/cmp"
	"gotest.tools/v3/assert"
	"testing"
)

type stringPointerWrapper *string

type Embedded struct{ int }

type testStruct struct {
	*Embedded
	Other *testStruct
	X     int
	SPtr  *string
	W     stringPointerWrapper
	Any   any
}

// type from the example
type Foo struct {
	SPtr *string
	S    string
}

func TestMake_Example(t *testing.T) {
	helloString := "hello"

	// true
	assert.Assert(t, cmp.Equal(Make[*string](t), &helloString, Ignore()))

	placeholder := Make[*string](t)
	anotherRef := &(*placeholder) // nolint

	// true, any reference to the allocated object is a placeholder
	assert.Assert(t, cmp.Equal(anotherRef, &helloString, Ignore()))

	// false, the allocated object itself is not a placeholder
	assert.Assert(t, !cmp.Equal(*placeholder, "hello", Ignore()))

	// true, it works with struct fields and embedded types as well!
	assert.Assert(t, cmp.Equal(Foo{Make[*string](t), "world"}, Foo{&helloString, "world"}, Ignore()))

	// false, non-placeholder fields differ
	assert.Assert(t, !cmp.Equal(Foo{Make[*string](t), "earthlings"}, Foo{&helloString, "world"}, Ignore()))
}

func TestMake(tt *testing.T) {
	helloString := "hello"
	worldString := "world"
	fiveInt := 5

	testCases := map[string]func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool){
		"zero values": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{}
			arg2 = testStruct{}
			equal = true
			return
		},

		"simple placeholder": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{X: 17, SPtr: Make[*string](t), W: &worldString}
			arg2 = testStruct{X: 17, SPtr: &helloString, W: &worldString}
			equal = true
			return
		},

		"simple placeholder against nil": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{X: 17, SPtr: Make[*string](t), W: &worldString}
			arg2 = testStruct{X: 17, W: &worldString}
			equal = true
			return
		},

		"embedded reference placeholder": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{Embedded: Make[*Embedded](t), X: 17, W: &worldString}
			arg2 = testStruct{Embedded: &Embedded{42}, X: 17, W: &worldString}
			equal = true
			return
		},

		"embedded reference placeholder against nil": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{Embedded: Make[*Embedded](t), X: 17, W: &worldString}
			arg2 = testStruct{X: 17, W: &worldString}
			equal = true
			return
		},

		"wrapped placeholder": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{X: 17, SPtr: &helloString, W: Make[stringPointerWrapper](t)}
			arg2 = testStruct{X: 17, SPtr: &helloString, W: &worldString}
			equal = true
			return
		},

		"wrapped placeholder against nil": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{X: 17, SPtr: &helloString, W: Make[stringPointerWrapper](t)}
			arg2 = testStruct{X: 17, SPtr: &helloString}
			equal = true
			return
		},

		"another reference to the created object": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			placeholder := Make[*string](t)
			anotherRef := &(*placeholder) // nolint

			if &placeholder == &anotherRef {
				panic("these are supposed to be two different pointers to the same object")
			}

			arg1 = testStruct{X: 17, SPtr: anotherRef, W: &worldString}
			arg2 = testStruct{X: 17, SPtr: &helloString, W: &worldString}
			equal = true
			return
		},

		"placeholder in a referenced struct": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{X: 17, SPtr: &helloString, W: Make[stringPointerWrapper](t)}
			arg1.Other = &testStruct{X: 42, SPtr: Make[*string](t)}

			arg2 = testStruct{X: 17, SPtr: &helloString, W: &worldString}
			arg2.Other = &testStruct{X: 42, SPtr: &helloString}

			equal = true
			return
		},

		"placeholder against recursive reference": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			s := "world"

			arg1 = testStruct{X: 17, SPtr: &s}
			arg1.Other = Make[*testStruct](t)

			arg2 = testStruct{X: 17, SPtr: &s}
			arg2.Other = &arg2

			equal = true
			return
		},

		"difference in a non-placeholder field": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			s := "world"
			arg1 = testStruct{X: 17, SPtr: Make[*string](t), W: Make[stringPointerWrapper](t)}
			arg2 = testStruct{X: 42, SPtr: &s, W: &s}
			equal = false
			return
		},

		"placeholder against a pointer to a different type": func(t *testing.T) (arg1 testStruct, arg2 testStruct, equal bool) {
			arg1 = testStruct{X: 17, SPtr: &helloString, Any: Make[*string](t)}
			arg2 = testStruct{X: 17, SPtr: &helloString, Any: &fiveInt}
			equal = true
			return
		},
	}

	for testName, tc := range testCases {
		tt.Run(testName, func(t *testing.T) {
			arg1, arg2, expectedEqual := tc(t)
			equal := cmp.Equal(arg1, arg2, Ignore())
			if expectedEqual && !equal {
				t.Errorf("structs are supposed to be considered equal, but are considered differet. diff: %v",
					cmp.Diff(arg1, arg2, Ignore()))
			} else if !expectedEqual && equal {
				t.Error("structs are supposed to be considered different, but are considered equal")
			}
		})
	}
}
