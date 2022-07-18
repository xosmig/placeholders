package matchers

import (
	"github.com/budougumi0617/cmpmock"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/xosmig/placeholders"
)

// DiffEq invokes cmpmock.DiffEq with an extra placeholders.Ignore() option.
func DiffEq(x any, opts ...cmp.Option) gomock.Matcher {
	return cmpmock.DiffEq(x, placeholders.Ignore(), cmp.Options(opts))
}
