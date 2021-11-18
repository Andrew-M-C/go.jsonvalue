package jsonvalue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIterFloat(t *testing.T) {
	test(t, "other parseResult conditions", testUnmarshalFloatErrors)
}

func testUnmarshalFloatErrors(t *testing.T) {
	Convey("overflow", func() {
		_, err := UnmarshalString(`-9223372036854775809`)
		So(err, ShouldBeError)
		_, err = UnmarshalString(`18446744073709551616`)
		So(err, ShouldBeError)
		_, err = UnmarshalString(`-9999999999999999999`)
		So(err, ShouldBeError)
		_, err = UnmarshalString(`9999999999999999999999999999999999999999999999999999999999999999999`)
		So(err, ShouldBeError)
		_, err = UnmarshalString(`99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999.9999999999999999999999999999999999999999999999999999999999999999999`)
		So(err, ShouldBeError)
	})

	Convey("stateStart", func() {
		it := &iter{b: []byte("E")}
		stm := newFloatStateMachine(it, 0)
		err := stm.next()
		So(err, ShouldBeError)
	})

	Convey("stateLeadingZero", func() {
		v, err := UnmarshalString(`0`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, 0)

		_, err = UnmarshalString(`01`)
		So(err, ShouldBeError)

		v, err = UnmarshalString(`0.0`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, 0)
		So(v.IsFloat(), ShouldBeTrue)

		v, err = UnmarshalString(`0E0`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, 0)
		So(v.IsFloat(), ShouldBeTrue)
	})

	Convey("stateLeadingDigit", func() {
		v, err := UnmarshalString(`1`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, 1)

		v, err = UnmarshalString(`1E1`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, 10)
		So(v.IsFloat(), ShouldBeTrue)

		v, err = UnmarshalString(`1E+1`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, 10)
		So(v.IsFloat(), ShouldBeTrue)

		_, err = UnmarshalString(`1Ee`)
		So(err, ShouldBeError)

		_, err = UnmarshalString(`1-`)
		So(err, ShouldBeError)
	})

	Convey("stateLeadingNegative", func() {
		v, err := UnmarshalString(`-1`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, -1)

		_, err = UnmarshalString(`-`)
		So(err, ShouldBeError)

		_, err = UnmarshalString(`-.`)
		So(err, ShouldBeError)

		v, err = UnmarshalString(`-0.25`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Float64(), ShouldEqual, -0.25)
	})

	Convey("stateIntegerDigit", func() {
		v, err := UnmarshalString(`10E-1`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Int(), ShouldEqual, 1)
		So(v.IsFloat(), ShouldBeTrue)

		_, err = UnmarshalString(`10-`)
		So(err, ShouldBeError)
	})

	Convey("stateExponent", func() {
		_, err := UnmarshalString(`1E`)
		So(err, ShouldBeError)

		_, err = UnmarshalString(`1E+`)
		So(err, ShouldBeError)

		_, err = UnmarshalString(`1e-`)
		So(err, ShouldBeError)
	})

	Convey("stateExponentSign", func() {
		_, err := UnmarshalString(`1E--`)
		So(err, ShouldBeError)
	})

	Convey("stateFractionDigit", func() {
		_, err := UnmarshalString(`1.1+`)
		So(err, ShouldBeError)
	})

	Convey("stateExponentDigit", func() {
		v, err := UnmarshalString(`-1e15`)
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.Float64(), ShouldEqual, -1e15)

		_, err = UnmarshalString(`1e2e`)
		So(err, ShouldBeError)
	})
}
