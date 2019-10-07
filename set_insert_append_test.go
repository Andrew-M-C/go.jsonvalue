package jsonvalue

import (
	"testing"
)

func TestSet(t *testing.T) {
	o := NewObject()
	child := NewString("Hello, world!")
	_, err := o.Set(child).At("data", "message", -1, "hello")
	if err != nil {
		t.Errorf("test Set failed: %v", err)
		return
	}

	b, _ := o.Marshal()
	t.Logf("after setting: %v", string(b))
	if string(b) != `{"data":{"message":[{"hello":"Hello, world!"}]}}` {
		t.Errorf("test Set() failed")
	}
	return
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

func TestSetStringNullBool(t *testing.T) {
	a := NewArray()
	expected := `[123456,"hello","world",1234.123456789,true,["12345"],null]`
	a.AppendString("world").InTheBeginning()
	a.AppendFloat64(1234.123456789, 9).InTheEnd()
	a.InsertBool(true).After(-1)
	a.AppendNull().InTheEnd()
	a.InsertInt(123456).Before(0)
	a.InsertString("hello").After(0)
	a.InsertArray().After(-2)
	a.AppendString("12345").InTheEnd(-2)

	s, _ := a.MarshalString()
	t.Logf("after SetXxx(): %v", s)
	if s != expected {
		t.Errorf("series SetXxx failed")
		return
	}
}
