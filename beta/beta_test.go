package beta

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So

	eq = convey.ShouldEqual

	isTrue  = convey.ShouldBeTrue
	isFalse = convey.ShouldBeFalse
)

func TestBeta(t *testing.T) {
	cv("test Contains()", t, func() { testContains(t) })
}
