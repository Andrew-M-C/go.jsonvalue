package jsonvalue

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestSet(t *testing.T) {
	o := NewObject()
	child := NewString("Hello, world!")
	_, err := o.Set(child).At("data", "message", 0, "hello")
	if err != nil {
		t.Errorf("test Set failed: %v", err)
		return
	}

	b, _ := o.Marshal()
	t.Logf("after setting: %v", string(b))
	if string(b) != `{"data":{"message":[{"hello":"Hello, world!"}]}}` {
		t.Errorf("test Set() failed")
	}
}

func TestSetInteger(t *testing.T) {
	var o *V
	var err error
	var s string
	const integer = 123456
	expected := `{"data":{"integer":123456}}`

	// SetInt()
	o = NewObject()
	_, err = o.SetInt(integer).At("data", "integer")
	if err != nil {
		t.Errorf("SetInt failed: %v", err)
		return
	}

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt:    %v", s)
	if s != expected {
		t.Errorf("test SetInt() failed")
		return
	}

	// SetInt32()
	o = NewObject()
	_, err = o.SetInt32(integer).At("data", "integer")
	if err != nil {
		t.Errorf("SetInt32 failed: %v", err)
		return
	}

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt32:  %v", s)
	if s != expected {
		t.Errorf("test SetInt32() failed")
		return
	}

	// SetInt64()
	o = NewObject()
	_, err = o.SetInt64(integer).At("data", "integer")
	if err != nil {
		t.Errorf("SetInt64 failed: %v", err)
		return
	}

	s, _ = o.MarshalString()
	t.Logf("\tafter SetInt64:  %v", s)
	if s != expected {
		t.Errorf("test SetInt64() failed")
		return
	}

	// SetUint()
	o = NewObject()
	_, err = o.SetUint(integer).At("data", "integer")
	if err != nil {
		t.Errorf("SetUint failed: %v", err)
		return
	}

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint:   %v", s)
	if s != expected {
		t.Errorf("test SetUint() failed")
		return
	}

	// SetUint64()
	o = NewObject()
	_, err = o.SetUint64(integer).At("data", "integer")
	if err != nil {
		t.Errorf("SetUint64 failed: %v", err)
		return
	}

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint64: %v", s)
	if s != expected {
		t.Errorf("test SetUint64() failed")
		return
	}

	// SetUint32()
	o = NewObject()
	_, err = o.SetUint32(integer).At("data", "integer")
	if err != nil {
		t.Errorf("SetUint32 failed: %v", err)
		return
	}

	s, _ = o.MarshalString()
	t.Logf("\tafter SetUint32: %v", s)
	if s != expected {
		t.Errorf("test SetUint32() failed")
		return
	}
}

func TestSetMisc(t *testing.T) {
	var err error
	var topic string
	var v *V
	checkErr := func() {
		if err != nil {
			t.Errorf("%s error: %v", topic, err)
			return
		}
	}
	check := func(b bool) {
		if false == b {
			t.Errorf("%s failed", topic)
			return
		}
	}

	v = NewObject()
	v.SetObject().At("data")

	topic = "SetBool(true)"
	v.SetBool(true).At("data", "true")
	b, err := v.GetBool("data", "true")
	checkErr()
	check(b == true)

	topic = "SetBool(false)"
	v.SetBool(false).At("data", "false")
	b, err = v.GetBool("data", "false")
	checkErr()
	check(b == false)

	topic = "SetFloat64"
	v.SetFloat64(1234.12345678, 8).At("data", "float64")
	f, err := v.Get("data", "float64")
	checkErr()
	check(f.String() == "1234.12345678")

	topic = "SetFloat32"
	v.SetFloat32(1234.123, 4).At("data", "float32")
	f, err = v.Get("data", "float32")
	checkErr()
	check(f.String() == "1234.1230")

	topic = "SetObject"
	v.SetObject().At("data", "object")
	v.SetString("hello").At("data", "object", "message")
	o, err := v.Get("data", "object")
	checkErr()
	check(o.IsObject() && o.Len() == 1)

	topic = "SetArray"
	v.SetArray().At("data", "array")
	v.AppendNull().InTheEnd("data", "array")
	a, err := v.Get("data", "array")
	checkErr()
	check(a.IsArray() && a.Len() == 1)

	topic = "GetBytes"
	s := "1234567890"
	data, _ := hex.DecodeString(s)
	v.SetString(s).At("string")
	v.SetBytes(data).At("bytes")
	dataRead, err := v.GetBytes("bytes")
	checkErr()
	t.Logf("set data: %s", hex.EncodeToString(data))
	t.Logf("Got data: %s", hex.EncodeToString(dataRead))
	check(bytes.Equal(data, dataRead))
	_, err = a.GetBytes("string")
	check(err != nil)

	topic = "Bytes"
	child, _ := v.Get("string")
	t.Logf("Get: %v", child)
	check(len(child.Bytes()) == 0)
	child, _ = v.Get("data")
	t.Logf("Get: %v", child)
	check(len(child.Bytes()) == 0)
	child, _ = v.Get("bytes")
	check(bytes.Equal(data, child.Bytes()))

	topic = "SetString in array of a object"
	a = NewArray()
	a.AppendObject().InTheBeginning()
	_, err = a.SetString("hello").At(0)
	checkErr()
	s, err = a.GetString(0)
	checkErr()
	check(s == "hello")

	topic = "Set(nil)"
	v.Set(nil).At("data", "nil")
	err = v.GetNull("data", "nil")
	checkErr()

	topic = "Complex Set()"
	a = NewArray()
	_, err = a.SetArray().At(0, 0, 0)
	checkErr()
	_, err = a.SetNull().At(0, 0, 1)
	checkErr()
	s, _ = a.MarshalString()
	check(s == "[[[[],null]]]")

	_, err = a.SetBool(true).At(0, 0, -1)
	checkErr()
	check(a.MustMarshalString() == "[[[[],true]]]")
}

func TestSetError(t *testing.T) {
	var checkCount int
	shouldError := func(err error) {
		checkCount++
		if err == nil {
			t.Errorf("%02d - error expected but not caught", checkCount)
			return
		}
		t.Logf("expected error string: %v", err)
	}

	{
		raw := `"`
		_, err := UnmarshalString(raw)
		shouldError(err)
	}

	{
		v := NewObject()
		_, err := v.SetString("hello").At(true)
		shouldError(err)
		_, err = v.SetString("hello").At(true, "message")
		shouldError(err)
		_, err = v.SetString("hello").At("message", true)
		shouldError(err)
		_, err = v.SetString("hello").At("message", "message", true)
		shouldError(err)
	}

	{
		v := NewObject()
		c := &V{}
		_, err := v.Set(c).At("uninitialized")
		shouldError(err)
		v.SetString("hello").At("object", "message")
		_, err = v.SetNull().At("object", "message", "null")
		shouldError(err)
		t.Logf("v: %s", v.MustMarshalString())
	}

	{
		v := &V{}
		c := NewObject()
		_, err := v.Set(c).At("uninitialized")
		shouldError(err)
	}

	{
		v := NewString("string")
		_, err := v.SetString("hello").At("message")
		shouldError(err)
		_, err = v.SetString("hello").At("object", "message")
		shouldError(err)
	}

	{
		v := NewArray()
		_, err := v.SetNull().At("0")
		shouldError(err)
		_, err = v.SetNull().At(1)
		shouldError(err)
	}

	{
		v := NewArray()
		v.AppendArray().InTheBeginning()
		v.AppendArray().InTheBeginning(0)
		v.AppendObject().InTheEnd(0)
		_, err := v.SetNull().At(0, true)
		shouldError(err)
		_, err = v.SetNull().At(0, 0, true)
		shouldError(err)
		_, err = v.SetNull().At(0, true, 0)
		shouldError(err)
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
		shouldError(err)
		_, err = v.SetNull().At(-10)
		shouldError(err)
	}

	{
		v := NewArray()
		_, err := v.SetNull().At(0, 1)
		shouldError(err)
		_, err = v.SetNull().At(0, true)
		shouldError(err)
	}

	{
		v := NewObject()
		_, err := v.SetNull().At("array", 1)
		shouldError(err)
		_, err = v.SetNull().At("array", true)
		shouldError(err)
	}
}
