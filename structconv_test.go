package jsonvalue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStructConv(t *testing.T) {
	test(t, "export to string", testExportString)
	test(t, "export to int", testExportInt)
	test(t, "export to float", testExportFloat)
	test(t, "export to bool", testExportBool)
}

func testExportString(t *testing.T) {
	const S = "Hello, jsonvalue!"
	v := NewString(S)

	str := ""
	err := v.Export(str)
	So(err, ShouldBeError)

	err = v.Export(&str)
	So(err, ShouldBeNil)

	So(str, ShouldEqual, S)

	bol := true
	err = v.Export(&bol)
	So(err, ShouldBeError)
}

func testExportInt(t *testing.T) {
	const positive = 123454321
	const negative = -987656789

	n1 := NewInt(positive)

	var i int
	var u uint
	var i32 int32
	var u32 uint32

	err := n1.Export(&i)
	So(err, ShouldBeNil)
	So(i, ShouldEqual, positive)

	err = n1.Export(&u)
	So(err, ShouldBeNil)
	So(u, ShouldEqual, positive)

	err = n1.Export(&i32)
	So(err, ShouldBeNil)
	So(i32, ShouldEqual, positive)

	err = n1.Export(&u32)
	So(err, ShouldBeNil)
	So(u32, ShouldEqual, positive)

	// --------

	n2 := NewInt(negative)

	err = n2.Export(&i)
	So(err, ShouldBeNil)
	So(i, ShouldEqual, negative)

	err = n2.Export(&i32)
	So(err, ShouldBeNil)
	So(i32, ShouldEqual, negative)

	// --------

	bol := true
	err = n1.Export(&bol)
	So(err, ShouldBeError)
}

func testExportFloat(t *testing.T) {
	const F = 12345.4321

	n := NewFloat64(F, 4)

	var f32 float32
	var f64 float64

	err := n.Export(&f32)
	So(err, ShouldBeNil)
	So(f32, ShouldEqual, F)

	err = n.Export(&f64)
	So(err, ShouldBeNil)
	So(f64, ShouldEqual, F)

	// --------

	bol := true
	err = n.Export(&bol)
	So(err, ShouldBeError)
}

func testExportBool(t *testing.T) {
	v := NewBool(true)
	b := false

	err := v.Export(b)
	So(err, ShouldBeError)

	err = v.Export(&b)
	So(err, ShouldBeNil)

	So(b, ShouldBeTrue)

	str := ""
	err = v.Export(&str)
	So(err, ShouldBeError)
}
