package beta_test

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So

	eq = convey.ShouldEqual

	resemble = convey.ShouldResemble

	isTrue  = convey.ShouldBeTrue
	isFalse = convey.ShouldBeFalse
)

type any = interface{}

func TestBeta(t *testing.T) {
	cv("test Contains()", t, func() { testContains(t) })
	cv("test Merge()", t, func() { testMerge(t) })
}
