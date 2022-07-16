package gomock

import (
	"github.com/golang/mock/gomock"
	"github.com/xosmig/placeholders"
	"github.com/xosmig/placeholders/examples/gomock/foo"
	"github.com/xosmig/placeholders/examples/gomock/foo/mock"
	"github.com/xosmig/placeholders/matchers"
	"gotest.tools/v3/assert"
	"testing"
)

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
