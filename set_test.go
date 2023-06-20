package jsonvalue

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func testSet(t *testing.T) {
	cv("general set", func() { testJsonvalue_Set(t) })
	cv("set integer", func() { testSetInteger(t) })
	cv("misc set", func() { testSetMisc(t) })
	cv("set errors", func() { testSetError(t) })
	cv("must set", func() { testSet_Must(t) })
}

func testJsonvalue_Set(t *testing.T) {
	o := NewObject()
	child := NewString("Hello, world!")

	_, err := o.Set(child).At("data", "message", 0, "hello")
	so(err, isNil)

	b, _ := o.Marshal()
	t.Logf("after setting: %v", string(b))
	so(string(b), eq, `{"data":{"message":[{"hello":"Hello, world!"}]}}`)
}

func testSetInteger(t *testing.T) {
	var o *V
	var err error
	var s string
	const integer = 123456
	expected := `{"data":{"integer":123456}}`

	// SetInt()
	o = NewObject()
	o.MustSetInt(integer).At("data", "integer")
	so(o.MustGet("data", "integer").Int(), eq, integer)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt:    %v", s)
	so(s, eq, expected)

	// SetInt32()
	o = NewObject()
	_, err = o.SetInt32(integer).At("data", "integer")
	so(err, isNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt32:  %v", s)
	so(s, eq, s)

	// SetInt64()
	o = NewObject()
	_, err = o.SetInt64(integer).At("data", "integer")
	so(err, isNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt64:  %v", s)
	so(s, eq, s)

	// SetUint()
	o = NewObject()
	_, err = o.SetUint(integer).At("data", "integer")
	so(err, isNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint:   %v", s)
	so(s, eq, expected)

	// SetUint64()
	o = NewObject()
	_, err = o.SetUint64(integer).At("data", "integer")
	so(err, isNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint64: %v", s)
	so(s, eq, expected)

	// SetUint32()
	o = NewObject()
	_, err = o.SetUint32(integer).At("data", "integer")
	so(err, isNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint32: %v", s)
	so(s, eq, expected)
}

func testSetMisc(t *testing.T) {
	var err error

	v := NewObject()
	v.MustSetObject().At("data")

	v.MustSetBool(true).At("data", "true")
	b, err := v.GetBool("data", "true")
	so(err, isNil)
	so(b, isTrue)

	v.MustSetBool(false).At("data", "false")
	b, err = v.GetBool("data", "false")
	so(err, isNil)
	so(b, isFalse)

	v.MustSetFloat64(1234.12345678).At("data", "float64")
	f, err := v.Get("data", "float64")
	so(err, isNil)
	so(f.String(), eq, "1234.12345678")

	v.MustSet(NewFloat32f(1234.123, 'f', 4)).At("data", "float32")
	f, err = v.Get("data", "float32")
	so(err, isNil)
	so(f.String(), eq, "1234.1230")

	v.MustSetFloat32(1234.123).At("data", "float32")
	f, err = v.Get("data", "float32")
	so(err, isNil)
	so(f.String(), eq, "1234.123")

	v.MustSetObject().At("data", "object")
	v.MustSetString("hello").At("data", "object", "message")
	o, err := v.Get("data", "object")
	so(err, isNil)
	so(o.IsObject(), isTrue)
	so(o.Len(), eq, 1)

	v.MustSetArray().At("data", "array")
	v.MustAppendNull().InTheEnd("data", "array")
	a, err := v.Get("data", "array")
	so(err, isNil)
	so(a.IsArray(), isTrue)
	so(a.Len(), eq, 1)

	s := "1234567890"
	data, _ := hex.DecodeString(s)
	v.MustSetString(s).At("string")
	v.MustSetBytes(data).At("bytes")
	dataRead, err := v.GetBytes("bytes")
	so(err, isNil)

	t.Logf("set data: %s", hex.EncodeToString(data))
	t.Logf("Got data: %s", hex.EncodeToString(dataRead))
	so(bytes.Equal(data, dataRead), isTrue)

	child, _ := v.Get("string")
	t.Logf("Get: %v", child)

	so(len(child.Bytes()), isZero)
	child, _ = v.Get("data")
	t.Logf("Get: %v", child)
	so(len(child.Bytes()), isZero)
	child, _ = v.Get("bytes")
	so(bytes.Equal(data, child.Bytes()), isTrue)

	a = NewArray()
	a.MustAppendObject().InTheBeginning()
	_, err = a.SetString("hello").At(0)
	so(err, isNil)
	s, err = a.GetString(0)
	so(err, isNil)
	so(s, eq, "hello")

	// Set(nil)
	v.MustSet(nil).At("data", "nil")
	err = v.GetNull("data", "nil")
	so(err, isNil)

	// Set(any type)
	v = NewObject()
	_, err = v.Set(make(chan struct{})).At("channel")
	so(err, notNil)
	child, err = v.Set([]any{1, "2", 3.25, -4, false, nil}).At("arr")
	so(child.IsArray(), isTrue)
	so(err, isNil)
	so(v.MustMarshalString(), eq, `{"arr":[1,"2",3.25,-4,false,null]}`)

	// Complex Set()
	a = NewArray()
	_, err = a.SetArray().At(0, 0, 0)
	so(err, isNil)
	_, err = a.SetNull().At(0, 0, 1)
	so(err, isNil)
	s, _ = a.MarshalString()
	so(s, eq, "[[[[],null]]]")

	_, err = a.SetBool(true).At(0, 0, -1)
	so(err, isNil)
	so(a.MustMarshalString(), eq, "[[[[],true]]]")
}

func testSetError(t *testing.T) {

	{
		raw := `"`
		_, err := UnmarshalString(raw)
		so(err, isErr)
	}

	{
		v := NewObject()
		_, err := v.SetString("hello").At(true)
		so(err, isErr)
		_, err = v.SetString("hello").At(true, "message")
		so(err, isErr)
		_, err = v.SetString("hello").At("message", true)
		so(err, isErr)
		_, err = v.SetString("hello").At("message", "message", true)
		so(err, isErr)
	}

	{
		v := NewObject()
		c := &V{}
		_, err := v.Set(c).At("uninitialized")
		so(err, isErr)
		v.MustSetString("hello").At("object", "message")
		_, err = v.SetNull().At("object", "message", "null")
		so(err, isErr)
		t.Logf("v: %s", v.MustMarshalString())
	}

	{
		v := &V{}
		c := NewObject()
		_, err := v.Set(c).At("uninitialized")
		so(err, isErr)
	}

	{
		v := NewString("string")
		_, err := v.SetString("hello").At("message")
		so(err, isErr)
		_, err = v.SetString("hello").At("object", "message")
		so(err, isErr)
	}

	{
		v := NewArray()
		_, err := v.SetNull().At("0")
		so(err, isErr)
		_, err = v.SetNull().At(1)
		so(err, isErr)
	}

	{
		v := NewArray()
		v.MustAppendArray().InTheBeginning()
		v.MustAppendArray().InTheBeginning(0)
		v.MustAppendObject().InTheEnd(0)
		_, err := v.SetNull().At(0, true)
		so(err, isErr)
		_, err = v.SetNull().At(0, 0, true)
		so(err, isErr)
		_, err = v.SetNull().At(0, true, 0)
		so(err, isErr)
	}

	{
		v := NewArray()
		v.MustSetNull().At(0)
		v.MustSetNull().At(1)
		if v.MustMarshalString() != `[null,null]` {
			t.Errorf("unexpected object: %v", v.MustMarshalString())
			return
		}
		_, err := v.SetNull().At(10)
		so(err, isErr)
		_, err = v.SetNull().At(-10)
		so(err, isErr)
	}

	{
		v := NewArray()
		_, err := v.SetNull().At(0, 1)
		so(err, isErr)
		_, err = v.SetNull().At(0, true)
		so(err, isErr)
	}

	{
		v := NewObject()
		_, err := v.SetNull().At("array", 1)
		so(err, isErr)
		_, err = v.SetNull().At("array", true)
		so(err, isErr)
	}
}

func testSet_Must(t *testing.T) {
	cv("SetBytes", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetBytes([]byte{1, 2}).At("bytes")
		b.MustSetBytes([]byte{1, 2}).At("bytes")
		so(a.Equal(b), isTrue)
	})

	cv("SetString", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetString("1").At("one")
		b.MustSetString("1").At("one")
		so(a.Equal(b), isTrue)
	})

	cv("SetInt", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetInt(-1).At("one")
		b.MustSetInt(-1).At("one")
		so(a.Equal(b), isTrue)
	})

	cv("SetUint", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetUint(1).At("one")
		b.MustSetUint(1).At("one")
		so(a.Equal(b), isTrue)
	})

	cv("SetInt32", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetInt32(-1).At("one")
		b.MustSetInt32(-1).At("one")
		so(a.Equal(b), isTrue)
	})

	cv("SetUint32", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetUint32(1).At("one")
		b.MustSetUint32(1).At("one")
		so(a.Equal(b), isTrue)
	})

	cv("SetFloat64", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetFloat64(1).At("one")
		b.MustSetFloat64(1).At("one")
		so(a.Equal(b), isTrue)
	})

	cv("SetFloat32", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetFloat32(1).At("one")
		b.MustSetFloat32(1).At("one")
		so(a.Equal(b), isTrue)
	})

	cv("SetObject", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetObject().At("obj")
		b.MustSetObject().At("obj")
		so(a.Equal(b), isTrue)
	})

	cv("SetArray", func() {
		a := NewObject()
		b := NewObject()
		_, _ = a.SetArray().At("arr")
		b.MustSetArray().At("arr")
		so(a.Equal(b), isTrue)
	})
}
