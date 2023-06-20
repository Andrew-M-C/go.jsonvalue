package jsonvalue

import (
	"encoding"
	"encoding/json"
	"reflect"
	"testing"
)

func testMarshalerUnmarshaler(t *testing.T) {
	cv("json.Marshaler & json.Unmarshaler", func() { testMarshalerUnmarshaler_JSON(t) })
	cv("encoding.BinaryMarshaler & encoding.BinaryUnmarshaler", func() { testMarshalerUnmarshaler_Binary(t) })
}

func testMarshalerUnmarshaler_JSON(t *testing.T) {
	cv("json.Marshaler", func() {
		v := NewObject()
		v.MustSet(1).At("one", "one")

		s := v.MustMarshalString(OptSetSequence())
		b, err := json.Marshal(v)
		so(err, isNil)
		so(string(b), eq, s)
	})

	cv("json.Unmarshaler", func() {
		raw := `[1, 2, "3", "4", null]`
		v := &V{}
		err := json.Unmarshal([]byte(raw), v)
		so(err, isNil)
		so(v.ValueType(), eq, Array)
		so(v.MustGet(0).String(), eq, "1")
		so(v.MustGet(0).ValueType(), eq, Number)
		so(v.MustGet(1).String(), eq, "2")
		so(v.MustGet(1).ValueType(), eq, Number)
		so(v.MustGet(2).String(), eq, "3")
		so(v.MustGet(2).ValueType(), eq, String)
		so(v.MustGet(3).String(), eq, "4")
		so(v.MustGet(3).ValueType(), eq, String)
		so(v.MustGet(4).String(), eq, "null")
		so(v.MustGet(4).ValueType(), eq, Null)
	})

	cv("json.Unmarshaler with error", func() {
		raw := `[1, 2, "3", "4", "null]` // lacking right \"
		v := &V{}
		err := json.Unmarshal([]byte(raw), v)
		so(err, isErr)

		u, ok := reflect.ValueOf(v).Interface().(json.Unmarshaler)
		so(ok, isTrue)
		err = u.UnmarshalJSON([]byte(raw))
		so(err, isErr)
	})
}

func testMarshalerUnmarshaler_Binary(t *testing.T) {
	cv("encoding.BinaryMarshaler", func() {
		v := NewObject()
		v.MustSet("sub-sub").At("obj", "obj")

		m, ok := reflect.ValueOf(v).Interface().(encoding.BinaryMarshaler)
		so(ok, isTrue)

		s := v.MustMarshalString(OptSetSequence())
		b, err := m.MarshalBinary()
		so(err, isNil)
		so(string(b), eq, s)
	})

	cv("encoding.BinaryUnmarshaler", func() {
		raw := `[1, 2, "3", "4", null]`
		v := &V{}

		u, ok := reflect.ValueOf(v).Interface().(encoding.BinaryUnmarshaler)
		so(ok, isTrue)

		err := u.UnmarshalBinary([]byte(raw))
		so(err, isNil)
		so(v.ValueType(), eq, Array)
		so(v.MustGet(0).String(), eq, "1")
		so(v.MustGet(0).ValueType(), eq, Number)
		so(v.MustGet(1).String(), eq, "2")
		so(v.MustGet(1).ValueType(), eq, Number)
		so(v.MustGet(2).String(), eq, "3")
		so(v.MustGet(2).ValueType(), eq, String)
		so(v.MustGet(3).String(), eq, "4")
		so(v.MustGet(3).ValueType(), eq, String)
		so(v.MustGet(4).String(), eq, "null")
		so(v.MustGet(4).ValueType(), eq, Null)
	})

	cv("encoding.BinaryUnmarshaler with error", func() {
		raw := `[1, 2, "3", "4", "null]` // lacking right \"
		v := &V{}
		u, ok := reflect.ValueOf(v).Interface().(encoding.BinaryUnmarshaler)
		so(ok, isTrue)
		err := u.UnmarshalBinary([]byte(raw))
		so(err, isErr)
	})
}
