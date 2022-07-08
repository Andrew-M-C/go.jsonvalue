package jsonvalue

import (
	"fmt"
	"testing"
)

func testEqual(t *testing.T) {
	cv("test simple types", func() { testEqualSimpleTypes(t) })
	cv("test number type", func() { testEqualNumbers(t) })
	cv("test object type", func() { testEqualObject(t) })
	cv("test array type", func() { testEqualArray(t) })
}

func testEqualSimpleTypes(t *testing.T) {
	cv("invalid type", func() {
		var v1, v2 *V
		so(v1.Equal(v2), isFalse)

		v1 = &V{}
		so(v1.Equal(v2), isFalse)

		v2 = &V{}
		so(v1.Equal(v2), isFalse)
	})

	cv("string", func() {
		v1 := New("Hello!")
		v2 := New("Hello")
		so(v1.Equal(v2), isFalse)

		v2 = New(v2.String() + "!")
		so(v1.Equal(v2), isTrue)
	})

	cv("boolean", func() {
		v1 := New(true)
		v2 := New(false)
		so(v1.Equal(v2), isFalse)

		v2 = New(true)
		so(v1.Equal(v2), isTrue)

		v1 = New(false)
		so(v1.Equal(v2), isFalse)

		v2 = New(false)
		so(v1.Equal(v2), isTrue)
	})

	cv("null", func() {
		v1 := New(nil)
		v2 := New(nil)
		so(v1.Equal(v2), isTrue)
	})

	cv("diff type", func() {
		v1 := NewObject()
		v2 := NewArray()
		so(v1.Equal(v2), isFalse)
	})
}

func testEqualNumbers(t *testing.T) {
	cv("float", func() {
		longFloat := fmt.Sprintf("%.30f", 20.20)
		v1, err := UnmarshalString(longFloat)
		so(err, isNil)

		v2 := New(20.20)

		t.Log("float:", longFloat)
		t.Log("v1:", v1)
		t.Log("v2:", v2)

		so(string(v1.srcByte), eq, longFloat)
		so(v1.String(), eq, longFloat)
		so(v1.String(), ne, v2.String())
		so(v1.Equal(v2), isFalse)

		v1 = New(20.20)
		so(v1.Equal(v2), isTrue)
	})

	cv("positive int", func() {
		v1 := MustUnmarshalString("10.0")
		v2 := MustUnmarshalString("10")
		so(v1.Equal(v2), isTrue)
	})
}

func testEqualObject(t *testing.T) {
	cv("general", func() {
		v1 := MustUnmarshalString(`{"obj":{},"arr":[]}`)
		v2 := MustUnmarshalString(`{"arr":[],"obj":{}}`)
		so(v1.Equal(v2), isTrue)

		v1 = MustUnmarshalString(`{"num":-1.0}`)
		v2 = MustUnmarshalString(`{"num":-1}`)
		so(v1.Equal(v2), isTrue)

		v1 = MustUnmarshalString(`{"obj":{"msg":"Hello, world!"}}`)
		v2 = MustUnmarshalString(`{"obj":{"msg":"Hello, world!"}}`)
		so(v1.Equal(v2), isTrue)

		v1 = MustUnmarshalString(`{"obj":{"msg":"Hello, world!"}}`)
		v2 = MustUnmarshalString(`{"obj":{"Msg":"Hello, world!"}}`)
		so(v1.Equal(v2), isFalse)

		v1 = MustUnmarshalString(`{"int":0,"str":""}`)
		v2 = MustUnmarshalString(`{"int":0}`)
		so(v1.Equal(v2), isFalse)
	})
}

func testEqualArray(t *testing.T) {
	cv("general", func() {
		v1 := MustUnmarshalString(`[1,2,3,4]`)
		v2 := MustUnmarshalString(`[1,2,3,4.0]`)
		so(v1.Equal(v2), isTrue)

		v1 = MustUnmarshalString(`[{"msg":"Hello, world"},2,3,4]`)
		v2 = MustUnmarshalString(`[{"msg":"Hello, world"},2,3,4]`)
		so(v1.Equal(v2), isTrue)

		v1 = MustUnmarshalString(`[{"msg":"Hello, world"},2,3,4]`)
		v2 = MustUnmarshalString(`[{"Msg":"Hello, world"},2,3,4]`)
		so(v1.Equal(v2), isFalse)

		v1 = MustUnmarshalString(`[2,{"msg":"Hello, world"},3,4]`)
		v2 = MustUnmarshalString(`[{"msg":"Hello, world"},2,3,4]`)
		so(v1.Equal(v2), isFalse)

		v1 = MustUnmarshalString(`[0,0]`)
		v2 = MustUnmarshalString(`[0]`)
		so(v1.Equal(v2), isFalse)
	})
}
