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

func testEqualSimpleTypes(*testing.T) {
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

func testEqualObject(*testing.T) {
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

func testEqualArray(*testing.T) {
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

func testGreaterThan(*testing.T) {
	cv("numbers", func() {
		i1234 := MustUnmarshalString(`1234`)
		f1234 := MustUnmarshalString(`1234.0`)
		i5678 := MustUnmarshalString(`5678`)
		f5678 := MustUnmarshalString(`5678.0`)
		so(i1234.IsFloat(), isFalse)
		so(f1234.IsFloat(), isTrue)
		so(i5678.IsFloat(), isFalse)
		so(f5678.IsFloat(), isTrue)
		so(i1234.GreaterThan(f1234), isFalse)
		so(f1234.GreaterThan(i1234), isFalse)
		so(i5678.GreaterThan(f1234), isTrue)
		so(f5678.GreaterThan(i1234), isTrue)
		so(i5678.GreaterThan(i1234), isTrue)
		so(f5678.GreaterThan(f1234), isTrue)
	})
	cv("other abnormal situations", func() {
		i5678 := MustUnmarshalString(`5678`)
		s1234 := MustUnmarshalString(`"1234"`)
		so(i5678.GreaterThan(s1234), isFalse)
		so(s1234.GreaterThan(i5678), isFalse)
	})
}

func testGreaterThanOrEqual(*testing.T) {
	cv("numbers", func() {
		i1234 := MustUnmarshalString(`1234`)
		ii1234 := MustUnmarshalString(`1234`)
		f1234 := MustUnmarshalString(`1234.0`)
		i5678 := MustUnmarshalString(`5678`)
		f5678 := MustUnmarshalString(`5678.0`)
		so(i1234.IsFloat(), isFalse)
		so(f1234.IsFloat(), isTrue)
		so(i5678.IsFloat(), isFalse)
		so(f5678.IsFloat(), isTrue)
		so(i1234.GreaterThanOrEqual(ii1234), isTrue)
		so(i1234.GreaterThanOrEqual(f1234), isTrue)
		so(f1234.GreaterThanOrEqual(i1234), isTrue)
		so(i5678.GreaterThanOrEqual(f1234), isTrue)
		so(f5678.GreaterThanOrEqual(i1234), isTrue)
		so(i5678.GreaterThanOrEqual(i1234), isTrue)
		so(f5678.GreaterThanOrEqual(f1234), isTrue)
	})
	cv("other abnormal situations", func() {
		i5678 := MustUnmarshalString(`5678`)
		s1234 := MustUnmarshalString(`"1234"`)
		so(i5678.GreaterThanOrEqual(s1234), isFalse)
		so(s1234.GreaterThanOrEqual(i5678), isFalse)
	})
}

func testLessThan(*testing.T) {
	cv("numbers", func() {
		i1234 := MustUnmarshalString(`1234`)
		f1234 := MustUnmarshalString(`1234.0`)
		i5678 := MustUnmarshalString(`5678`)
		f5678 := MustUnmarshalString(`5678.0`)
		so(i1234.IsFloat(), isFalse)
		so(f1234.IsFloat(), isTrue)
		so(i5678.IsFloat(), isFalse)
		so(f5678.IsFloat(), isTrue)
		so(i1234.LessThan(f1234), isFalse)
		so(f1234.LessThan(i1234), isFalse)
		so(i1234.LessThan(f5678), isTrue)
		so(f1234.LessThan(i5678), isTrue)
		so(i1234.LessThan(f5678), isTrue)
		so(f1234.LessThan(i5678), isTrue)
		so(i1234.LessThan(i5678), isTrue)
		so(f1234.LessThan(f5678), isTrue)
		so(f1234.LessThan(i5678), isTrue)
		so(i1234.LessThan(f5678), isTrue)
	})
	cv("other abnormal situations", func() {
		i5678 := MustUnmarshalString(`5678`)
		s1234 := MustUnmarshalString(`"1234"`)
		so(i5678.LessThan(s1234), isFalse)
		so(s1234.LessThan(i5678), isFalse)
	})
}

func testLessThanOrEqual(*testing.T) {
	cv("numbers", func() {
		i1234 := MustUnmarshalString(`1234`)
		f1234 := MustUnmarshalString(`1234.0`)
		i5678 := MustUnmarshalString(`5678`)
		f5678 := MustUnmarshalString(`5678.0`)
		so(i1234.IsFloat(), isFalse)
		so(f1234.IsFloat(), isTrue)
		so(i5678.IsFloat(), isFalse)
		so(f5678.IsFloat(), isTrue)
		so(i1234.LessThanOrEqual(f1234), isTrue)
		so(f1234.LessThanOrEqual(i1234), isTrue)
		so(i1234.LessThanOrEqual(f5678), isTrue)
		so(f1234.LessThanOrEqual(i5678), isTrue)
		so(i1234.LessThanOrEqual(f5678), isTrue)
		so(f1234.LessThanOrEqual(i5678), isTrue)
		so(i1234.LessThanOrEqual(i5678), isTrue)
		so(f1234.LessThanOrEqual(f5678), isTrue)
		so(f1234.LessThanOrEqual(i5678), isTrue)
		so(i1234.LessThanOrEqual(f5678), isTrue)
		so(i5678.LessThanOrEqual(f1234), isFalse)
		so(f5678.LessThanOrEqual(i1234), isFalse)
		so(i5678.LessThanOrEqual(i1234), isFalse)
		so(f5678.LessThanOrEqual(f1234), isFalse)
	})
	cv("other abnormal situations", func() {
		i5678 := MustUnmarshalString(`5678`)
		s1234 := MustUnmarshalString(`"1234"`)
		so(i5678.LessThanOrEqual(s1234), isFalse)
		so(s1234.LessThanOrEqual(i5678), isFalse)
	})
}
