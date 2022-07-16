package matchers

import (
	"github.com/budougumi0617/cmpmock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/xosmig/placeholders"
)

// DiffEq invokes cmpmock.DiffEq with an extra placeholders.Comparer() option.
// Note that, if another cmp.Comparer or cmp.Transformer option is provided, it
// will cause an ambiguity and gocmp will panic. Use a custom wrapper option in
// this case.
func DiffEq(x any, opts ...cmp.Option) gomock.Matcher {
	return cmpmock.DiffEq(x, append(opts, placeholders.Comparer())...)
}
