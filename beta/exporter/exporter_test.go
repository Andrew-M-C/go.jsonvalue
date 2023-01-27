package exporter

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So
	eq = convey.ShouldEqual
)

func TestExporter(t *testing.T) {
	internal.debugf = t.Logf

	cv("调试", t, func() { testDebugging(t) })
}

type simpleStruct struct {
	S string `json:"s"`
	I int    `json:"i"`
}

func testDebugging(t *testing.T) {
	st := simpleStruct{
		S: "Hello",
		I: 2023,
	}

	e, err := ParseExporter(st)
	so(err, eq, nil)

	t.Logf("Got: %v", e)

	v := e.Export(st)
	t.Log(v.MustMarshalString())
}
