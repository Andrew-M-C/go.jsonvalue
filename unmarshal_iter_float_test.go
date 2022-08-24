package jsonvalue

import (
	"testing"
)

func testIterFloat(t *testing.T) {
	cv("other parseResult conditions", func() { testUnmarshalFloatErrors(t) })
	cv("https://github.com/Andrew-M-C/go.jsonvalue/issues/8", func() { testIssue8(t) })
}

func testUnmarshalFloatErrors(t *testing.T) {
	cv("overflow", func() {
		_, err := UnmarshalString(`-9223372036854775809`)
		so(err, isErr)
		_, err = UnmarshalString(`18446744073709551616`)
		so(err, isErr)
		_, err = UnmarshalString(`-9999999999999999999`)
		so(err, isErr)
		_, err = UnmarshalString(`9999999999999999999999999999999999999999999999999999999999999999999`)
		so(err, isErr)
		_, err = UnmarshalString(`99999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999.9999999999999999999999999999999999999999999999999999999999999999999`)
		so(err, isErr)
	})

	cv("stateStart", func() {
		u := newUnmarshaler([]byte{'E'})
		_, _, _, err := u.parseNumber(0)
		so(err, isErr)
	})

	cv("stateLeadingZero", func() {
		v, err := UnmarshalString(`0`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 0)

		_, err = UnmarshalString(`01`)
		so(err, isErr)

		_, err = UnmarshalString(`00`)
		so(err, isErr)

		_, err = UnmarshalString(`+1`)
		so(err, isErr)

		_, err = UnmarshalString(`-00`)
		so(err, isErr)

		v, err = UnmarshalString(`0.0`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 0)
		so(v.IsFloat(), isTrue)

		v, err = UnmarshalString(`0.10`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 0)
		so(v.IsFloat(), isTrue)
		so(v.Float64(), eq, 0.1)

		v, err = UnmarshalString(`0E0`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 0)
		so(v.IsFloat(), isTrue)
	})

	cv("stateLeadingDigit", func() {
		v, err := UnmarshalString(`1`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 1)

		v, err = UnmarshalString(`1E1`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 10)
		so(v.IsFloat(), isTrue)

		v, err = UnmarshalString(`1E+1`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 10)
		so(v.IsFloat(), isTrue)

		_, err = UnmarshalString(`1Ee`)
		so(err, isErr)

		_, err = UnmarshalString(`1-`)
		so(err, isErr)
	})

	cv("stateLeadingNegative", func() {
		v, err := UnmarshalString(`-1`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, -1)

		_, err = UnmarshalString(`-`)
		so(err, isErr)

		_, err = UnmarshalString(`-.`)
		so(err, isErr)

		v, err = UnmarshalString(`-0.25`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Float64(), eq, -0.25)
	})

	cv("stateIntegerDigit", func() {
		v, err := UnmarshalString(`10E-1`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Int(), eq, 1)
		so(v.IsFloat(), isTrue)

		_, err = UnmarshalString(`10-`)
		so(err, isErr)
	})

	cv("stateExponent", func() {
		_, err := UnmarshalString(`1E`)
		so(err, isErr)

		_, err = UnmarshalString(`1E+`)
		so(err, isErr)

		_, err = UnmarshalString(`1e-`)
		so(err, isErr)
	})

	cv("stateExponentSign", func() {
		_, err := UnmarshalString(`1E--`)
		so(err, isErr)
	})

	cv("stateFractionDigit", func() {
		_, err := UnmarshalString(`1.1+`)
		so(err, isErr)
	})

	cv("stateExponentDigit", func() {
		v, err := UnmarshalString(`-1e15`)
		so(err, isNil)
		so(v, notNil)
		so(v.IsNumber(), isTrue)
		so(v.Float64(), eq, -1e15)

		_, err = UnmarshalString(`1e2e`)
		so(err, isErr)
	})
}

func testIssue8(t *testing.T) {
	strJson := []byte(`{"tunnels":[{"name":"command_line","uri":"/api/tunnels/command_line","public_url":"https://11111.ngrok.io","proto":"https","config":{"addr":"http://localhost:11111","inspect":true},"metrics":{"conns":{"count":1,"gauge":0,"rate1":5.456067032277228e-19,"rate5":0.0000016821504265361616,"rate15":0.00008846097772300972,"p50":8287268034,"p90":8287268034,"p95":8287268034,"p99":8287268034},"http":{"count":5,"rate1":2.5535363027836646e-18,"rate5":0.000008299538128664852,"rate15":0.0004403445395661658,"p50":427625,"p90":600127,"p95":600127,"p99":600127}}},{"name":"command_line (http)","uri":"/api/tunnels/command_line%20%28http%29","public_url":"http://11111.ngrok.io","proto":"http","config":{"addr":"http://localhost:11111","inspect":true},"metrics":{"conns":{"count":0,"gauge":0,"rate1":0,"rate5":0,"rate15":0,"p50":0,"p90":0,"p95":0,"p99":0},"http":{"count":0,"rate1":0,"rate5":0,"rate15":0,"p50":0,"p90":0,"p95":0,"p99":0}}}],"uri":"/api/tunnels"}`)
	j, err := Unmarshal(strJson)
	so(err, isNil)

	bb, _ := j.GetString("tunnels", 0, "proto")
	so(bb, eq, "https")

	cc, _ := j.GetString("tunnels", 0, "public_url")
	so(cc, eq, "https://11111.ngrok.io")

	v, err := j.Get("tunnels", 0, "metrics", "conns", "rate1")
	so(err, isNil)
	so(v.IsFloat(), isTrue)
	so(v.String(), eq, "5.456067032277228e-19")
}
