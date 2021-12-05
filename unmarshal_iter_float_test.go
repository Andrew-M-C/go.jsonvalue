package jsonvalue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIterFloat(t *testing.T) {
	test(t, "test floatStateMachine", testFloatStateMachine)
	test(t, "other parseResult conditions", testUnmarshalFloatErrors)
	test(t, "https://github.com/Andrew-M-C/go.jsonvalue/issues/8", testIssue8)
}

func testFloatStateMachine(t *testing.T) {
	Convey("basic", func() {
		it := iter{'0'}
		stm := newFloatStateMachine(it, 0)
		So(stm.offset(), ShouldBeZeroValue)

		stm = stm.withOffsetAddOne()
		So(stm.offset(), ShouldEqual, 1)
	})
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
		it := &iter{'E'}
		_, _, _, err := it.parseNumber(0)
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

func testIssue8(t *testing.T) {
	strJson := []byte(`{"tunnels":[{"name":"command_line","uri":"/api/tunnels/command_line","public_url":"https://11111.ngrok.io","proto":"https","config":{"addr":"http://localhost:11111","inspect":true},"metrics":{"conns":{"count":1,"gauge":0,"rate1":5.456067032277228e-19,"rate5":0.0000016821504265361616,"rate15":0.00008846097772300972,"p50":8287268034,"p90":8287268034,"p95":8287268034,"p99":8287268034},"http":{"count":5,"rate1":2.5535363027836646e-18,"rate5":0.000008299538128664852,"rate15":0.0004403445395661658,"p50":427625,"p90":600127,"p95":600127,"p99":600127}}},{"name":"command_line (http)","uri":"/api/tunnels/command_line%20%28http%29","public_url":"http://11111.ngrok.io","proto":"http","config":{"addr":"http://localhost:11111","inspect":true},"metrics":{"conns":{"count":0,"gauge":0,"rate1":0,"rate5":0,"rate15":0,"p50":0,"p90":0,"p95":0,"p99":0},"http":{"count":0,"rate1":0,"rate5":0,"rate15":0,"p50":0,"p90":0,"p95":0,"p99":0}}}],"uri":"/api/tunnels"}`)
	j, err := Unmarshal(strJson)
	So(err, ShouldBeNil)

	bb, _ := j.GetString("tunnels", 0, "proto")
	So(bb, ShouldEqual, "https")

	cc, _ := j.GetString("tunnels", 0, "public_url")
	So(cc, ShouldEqual, "https://11111.ngrok.io")

	v, err := j.Get("tunnels", 0, "metrics", "conns", "rate1")
	So(err, ShouldBeNil)
	So(v.IsFloat(), ShouldBeTrue)
	So(v.String(), ShouldEqual, "5.456067032277228e-19")
}
