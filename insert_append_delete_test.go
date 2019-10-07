package jsonvalue

import (
	"testing"
)

func TestInsertAppend(t *testing.T) {
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

func TestDelete(t *testing.T) {
	raw := `{"array":[1,2,3,4,5,6],"string":"string to be deleted","object":{"number":12345}}`
	o, err := UnmarshalString(raw)
	if err != nil {
		t.Errorf("UnmarshalString failed: %v", err)
		return
	}
	s, _ := o.MarshalString()
	t.Logf("parsed object: %v", s)

	err = o.Delete("object", "number")
	if err != nil {
		t.Errorf("Delete error: %v", err)
		return
	}

	sub, err := o.Get("object")
	if err != nil {
		t.Errorf("Get object failed: %v", err)
		return
	}

	s, _ = sub.MarshalString()
	if s != "{}" {
		t.Errorf("get sub mismatch: '%s'", s)
		return
	}

	err = o.Delete("object", "number")
	if err != ErrNotFound {
		t.Errorf("delete inexisted object, error should be raised (%v)", err)
		return
	}

	err = o.Delete("object")
	if err != nil {
		t.Errorf("Delete 'object' failed: %v", err)
		return
	}

	err = o.Delete("string")
	if err != nil {
		t.Errorf("Delete 'string' failed: %v", err)
		return
	}

	s, _ = o.MarshalString()
	if s != `{"array":[1,2,3,4,5,6]}` {
		t.Errorf("Deleted result not expected: '%s'", s)
		return
	}
}
