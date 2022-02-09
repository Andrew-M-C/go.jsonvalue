package jsonvalue

import (
	"bytes"
	"encoding/hex"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSet(t *testing.T) {
	test(t, "general set", testSet)
	test(t, "set integer", testSetInteger)
	test(t, "misc set", testSetMisc)
	test(t, "set errors", testSetError)
}

func testSet(t *testing.T) {
	o := NewObject()
	child := NewString("Hello, world!")

	_, err := o.Set(child).At("data", "message", 0, "hello")
	So(err, ShouldBeNil)

	b, _ := o.Marshal()
	t.Logf("after setting: %v", string(b))
	So(string(b), ShouldEqual, `{"data":{"message":[{"hello":"Hello, world!"}]}}`)
}

func testSetInteger(t *testing.T) {
	var o *V
	var err error
	var s string
	const integer = 123456
	expected := `{"data":{"integer":123456}}`

	// SetInt()
	o = NewObject()
	_, err = o.SetInt(integer).At("data", "integer")
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt:    %v", s)
	So(s, ShouldEqual, expected)

	// SetInt32()
	o = NewObject()
	_, err = o.SetInt32(integer).At("data", "integer")
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt32:  %v", s)
	So(s, ShouldEqual, s)

	// SetInt64()
	o = NewObject()
	_, err = o.SetInt64(integer).At("data", "integer")
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt64:  %v", s)
	So(s, ShouldEqual, s)

	// SetUint()
	o = NewObject()
	_, err = o.SetUint(integer).At("data", "integer")
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint:   %v", s)
	So(s, ShouldEqual, expected)

	// SetUint64()
	o = NewObject()
	_, err = o.SetUint64(integer).At("data", "integer")
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint64: %v", s)
	So(s, ShouldEqual, expected)

	// SetUint32()
	o = NewObject()
	_, err = o.SetUint32(integer).At("data", "integer")
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint32: %v", s)
	So(s, ShouldEqual, expected)
}

func testSetMisc(t *testing.T) {
	var err error

	v := NewObject()
	v.SetObject().At("data")

	v.SetBool(true).At("data", "true")
	b, err := v.GetBool("data", "true")
	So(err, ShouldBeNil)
	So(b, ShouldBeTrue)

	v.SetBool(false).At("data", "false")
	b, err = v.GetBool("data", "false")
	So(err, ShouldBeNil)
	So(b, ShouldBeFalse)

	v.SetFloat64(1234.12345678).At("data", "float64")
	f, err := v.Get("data", "float64")
	So(err, ShouldBeNil)
	So(f.String(), ShouldEqual, "1234.12345678")

	v.Set(NewFloat32f(1234.123, 'f', 4)).At("data", "float32")
	f, err = v.Get("data", "float32")
	So(err, ShouldBeNil)
	So(f.String(), ShouldEqual, "1234.1230")

	v.SetFloat32(1234.123).At("data", "float32")
	f, err = v.Get("data", "float32")
	So(err, ShouldBeNil)
	So(f.String(), ShouldEqual, "1234.123")

	v.SetObject().At("data", "object")
	v.SetString("hello").At("data", "object", "message")
	o, err := v.Get("data", "object")
	So(err, ShouldBeNil)
	So(o.IsObject(), ShouldBeTrue)
	So(o.Len(), ShouldEqual, 1)

	v.SetArray().At("data", "array")
	v.AppendNull().InTheEnd("data", "array")
	t.Log(v.MustMarshalString(OptDefaultStringSequence()))
	a, err := v.Get("data", "array")
	So(err, ShouldBeNil)
	So(a.IsArray(), ShouldBeTrue)
	So(a.Len(), ShouldEqual, 1)

	s := "1234567890"
	data, _ := hex.DecodeString(s)
	v.SetString(s).At("string")
	v.SetBytes(data).At("bytes")
	dataRead, err := v.GetBytes("bytes")
	So(err, ShouldBeNil)

	t.Logf("set data: %s", hex.EncodeToString(data))
	t.Logf("Got data: %s", hex.EncodeToString(dataRead))
	So(bytes.Equal(data, dataRead), ShouldBeTrue)

	child, _ := v.Get("string")
	t.Logf("Get: %v", child)

	So(len(child.Bytes()), ShouldBeZeroValue)
	child, _ = v.Get("data")
	t.Logf("Get: %v", child)
	So(len(child.Bytes()), ShouldBeZeroValue)
	child, _ = v.Get("bytes")
	So(bytes.Equal(data, child.Bytes()), ShouldBeTrue)

	a = NewArray()
	a.AppendObject().InTheBeginning()
	_, err = a.SetString("hello").At(0)
	So(err, ShouldBeNil)
	s, err = a.GetString(0)
	So(err, ShouldBeNil)
	So(s, ShouldEqual, "hello")

	// Set(nil)
	v.Set(nil).At("data", "nil")
	err = v.GetNull("data", "nil")
	So(err, ShouldBeNil)

	// Complex Set()
	a = NewArray()
	_, err = a.SetArray().At(0, 0, 0)
	So(err, ShouldBeNil)
	_, err = a.SetNull().At(0, 0, 1)
	So(err, ShouldBeNil)
	s, _ = a.MarshalString()
	So(s, ShouldEqual, "[[[[],null]]]")

	_, err = a.SetBool(true).At(0, 0, -1)
	So(err, ShouldBeNil)
	So(a.MustMarshalString(), ShouldEqual, "[[[[],true]]]")
}

func testSetError(t *testing.T) {

	{
		raw := `"`
		_, err := UnmarshalString(raw)
		So(err, ShouldBeError)
	}

	{
		v := NewObject()
		_, err := v.SetString("hello").At(true)
		So(err, ShouldBeError)
		_, err = v.SetString("hello").At(true, "message")
		So(err, ShouldBeError)
		_, err = v.SetString("hello").At("message", true)
		So(err, ShouldBeError)
		_, err = v.SetString("hello").At("message", "message", true)
		So(err, ShouldBeError)
	}

	{
		v := NewObject()
		c := &V{}
		_, err := v.Set(c).At("uninitialized")
		So(err, ShouldBeError)
		v.SetString("hello").At("object", "message")
		_, err = v.SetNull().At("object", "message", "null")
		So(err, ShouldBeError)
		t.Logf("v: %s", v.MustMarshalString())
	}

	{
		v := &V{}
		c := NewObject()
		_, err := v.Set(c).At("uninitialized")
		So(err, ShouldBeError)
	}

	{
		v := NewString("string")
		_, err := v.SetString("hello").At("message")
		So(err, ShouldBeError)
		_, err = v.SetString("hello").At("object", "message")
		So(err, ShouldBeError)
	}

	{
		v := NewArray()
		_, err := v.SetNull().At("0")
		So(err, ShouldBeError)
		_, err = v.SetNull().At(1)
		So(err, ShouldBeError)
	}

	{
		v := NewArray()
		v.AppendArray().InTheBeginning()
		v.AppendArray().InTheBeginning(0)
		v.AppendObject().InTheEnd(0)
		_, err := v.SetNull().At(0, true)
		So(err, ShouldBeError)
		_, err = v.SetNull().At(0, 0, true)
		So(err, ShouldBeError)
		_, err = v.SetNull().At(0, true, 0)
		So(err, ShouldBeError)
	}

	{
		v := NewArray()
		v.SetNull().At(0)
		v.SetNull().At(1)
		if v.MustMarshalString() != `[null,null]` {
			t.Errorf("unexpected object: %v", v.MustMarshalString())
			return
		}
		_, err := v.SetNull().At(10)
		So(err, ShouldBeError)
		_, err = v.SetNull().At(-10)
		So(err, ShouldBeError)
	}

	{
		v := NewArray()
		_, err := v.SetNull().At(0, 1)
		So(err, ShouldBeError)
		_, err = v.SetNull().At(0, true)
		So(err, ShouldBeError)
	}

	{
		v := NewObject()
		_, err := v.SetNull().At("array", 1)
		So(err, ShouldBeError)
		_, err = v.SetNull().At("array", true)
		So(err, ShouldBeError)
	}
}
