package beta

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So

	isTrue  = convey.ShouldBeTrue
	isFalse = convey.ShouldBeFalse
)

func TestBeta(t *testing.T) {
	cv("test HasSubset() and IsSubsetOf()", t, func() { testHasSubsetIsSubsetOf(t) })
}
