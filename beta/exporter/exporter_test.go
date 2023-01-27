package exporter

import (
	"testing"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So
	eq = convey.ShouldEqual
)

func TestExporter(t *testing.T) {
	internal.debugf = t.Logf

	cv("最基本的 struct", t, func() { testSimplestStruct(t) })
	cv("struct 嵌套自己", t, func() { testNestedStruct(t) })
}

type simplestStruct struct {
	S string `json:"s"`
	I int    `json:"i"`
}

func testSimplestStruct(t *testing.T) {
	st := simplestStruct{
		S: "Hello",
		I: 2023,
	}

	e, err := ParseExporter(st)
	so(err, eq, nil)

	t.Logf("Got: %v", e)

	v := e.Export(st)
	s := v.MustMarshalString(jsonvalue.OptSetSequence())
	t.Log(s)
	so(s, eq, `{"s":"Hello","i":2023}`)
}

type nestedStruct struct {
	ID string `json:"id"`

	SubWithEmpty *nestedStruct `json:"sub_with_empty"`
	Sub          *nestedStruct `json:"sub,omitempty"`
}

func testNestedStruct(t *testing.T) {
	st := nestedStruct{
		ID: "parent",
		Sub: &nestedStruct{
			ID: "child",
		},
	}

	e, err := ParseExporter(st)
	so(err, eq, nil)

	t.Log("Got:", e)

	v := e.Export(st)
	s := v.MustMarshalString(jsonvalue.OptSetSequence())
	t.Log(s)
	so(s, eq, `{"id":"parent","sub_with_empty":null,"sub":{"id":"child","sub_with_empty":null}}`)
}
