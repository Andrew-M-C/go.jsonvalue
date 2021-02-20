package jsonvalue

import (
	"testing"
)

func TestInsertAppend(t *testing.T) {
	a := NewArray()
	expected := `[123456,"hello","world",1234.123456789,true,["12345"],null,null]`
	a.AppendString("world").InTheBeginning()
	a.AppendFloat64(1234.123456789, 9).InTheEnd()
	a.InsertBool(true).After(-1)
	a.AppendNull().InTheEnd()
	a.InsertInt(123456).Before(0)
	a.InsertString("hello").After(0)
	a.InsertArray().After(-2)
	a.AppendString("12345").InTheEnd(-2)
	a.Append(nil).InTheEnd()

	s, _ := a.MarshalString()
	t.Logf("after SetXxx(): %v", s)
	if s != expected {
		t.Errorf("series SetXxx failed")
		return
	}
}

func TestDelete(t *testing.T) {
	raw := `{"array":[1,2,3,4,5,6],"string":"string to be deleted","object":{"number":12345},"Object":{}}`
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

	sub, err = o.Get("object")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	err = o.Delete("object") // delete another "object", actually "Object"
	if err != nil {
		t.Errorf("Delete 'object' failed: %v", err)
		return
	}

	_, err = o.Get("object")
	if err != ErrNotFound {
		t.Errorf("unexpected error: %v", err)
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

	err = o.Delete("array", 1)
	s, _ = o.MarshalString()
	if err != nil {
		t.Errorf("Delete array item 1 failed: %v", err)
		return
	}
	if s != `{"array":[1,3,4,5,6]}` {
		t.Errorf("Deleted result not expected: '%s'", s)
		return
	}
}

func TestMiscAppend(t *testing.T) {
	expected := `[true,-1,2,-3,4,-5,6,-7.7,8.8000,{},[[false],null]]`
	a := NewArray()
	a.AppendBool(true).InTheBeginning()
	a.AppendInt(-1).InTheEnd()
	a.AppendUint(2).InTheEnd()
	a.AppendInt32(-3).InTheEnd()
	a.AppendUint32(4).InTheEnd()
	a.AppendInt64(-5).InTheEnd()
	a.AppendUint64(6).InTheEnd()
	a.AppendFloat32(-7.7, 1).InTheEnd()
	a.AppendFloat64(8.8, 4).InTheEnd()
	a.AppendObject().InTheEnd()
	a.AppendArray().InTheEnd()
	a.AppendNull().InTheEnd(-1)
	a.AppendArray().InTheBeginning(-1)
	a.AppendBool(false).InTheBeginning(-1, 0)

	s, _ := a.MarshalString()
	if s != expected {
		t.Errorf("string not expected: '%s'", s)
		return
	}
}

func TestMiscInsert(t *testing.T) {
	var checkCount int
	var err error
	var c *V
	v := NewArray()
	expected := `[null,1,-2,3,-4,5,-6,7.7,-8.88888,true,false,null,null,{},[[null,-11,22]]]`

	checkErr := func(err error) {
		s, _ := v.MarshalString()
		t.Logf("marshaled: '%s'", s)

		if err != nil {
			t.Errorf("%02d - unexpected error: %v", checkCount, err)
		}
		return
	}
	checkCond := func(b bool) {
		if false {
			t.Errorf("%02d - check failed", checkCount)
		}
		checkCount++
		return
	}

	v.AppendNull().InTheBeginning()

	c, err = v.InsertUint(1).After(-1)
	checkErr(err)
	checkCond(c.Int() == 1)

	c, err = v.InsertInt(-2).After(-1)
	checkErr(err)
	checkCond(c.Int() == -2)

	c, err = v.InsertUint64(3).After(-1)
	checkErr(err)
	checkCond(c.Int() == 3)

	c, err = v.InsertInt64(-4).After(-1)
	checkErr(err)
	checkCond(c.Int() == -4)

	c, err = v.InsertUint32(5).After(-1)
	checkErr(err)
	checkCond(c.Int() == 5)

	c, err = v.InsertInt32(-6).After(-1)
	checkErr(err)
	checkCond(c.Int() == -6)

	c, err = v.InsertFloat32(7.7, -1).After(-1)
	checkErr(err)
	checkCond(c.Float64() == 7.7)

	c, err = v.InsertFloat64(-8.88888, 5).After(-1)
	checkErr(err)
	checkCond(c.Float64() == -8.8888)

	c, err = v.InsertBool(true).After(-1)
	checkErr(err)
	checkCond(c.Bool() == true)

	c, err = v.InsertBool(false).After(-1)
	checkErr(err)
	checkCond(c.Bool() == false && c.IsBoolean())

	c, err = v.Insert(nil).After(-1)
	checkErr(err)
	checkCond(c.IsNull())

	c, err = v.InsertNull().After(-1)
	checkErr(err)
	checkCond(c.IsNull())

	c, err = v.InsertObject().After(-1)
	checkErr(err)
	checkCond(c.IsObject())

	c, err = v.InsertArray().After(-1)
	checkErr(err)
	checkCond(c.IsArray())

	c, err = v.AppendArray().InTheBeginning(-1)
	checkErr(err)
	checkCond(c.IsArray())

	c, err = v.AppendInt(-11).InTheBeginning(-1, 0)
	checkErr(err)
	checkCond(c.Int() == -11)

	c, err = v.InsertUint(22).After(-1, 0, 0)
	checkErr(err)
	checkCond(c.Int() == 22)

	c, err = v.InsertNull().Before(-1, 0, 0)
	checkErr(err)
	checkCond(c.IsNull())

	s, _ := v.MarshalString()
	if s != expected {
		t.Errorf("marshaled string not expected: '%s'", s)
		return
	}
}

func TestMiscInsertError(t *testing.T) {
	var topic string
	var checkCount int
	shouldError := func(err error) {
		defer func() {
			checkCount++
		}()
		if err == nil {
			t.Errorf("%02d - check '%s' - error expected but not caught", checkCount, topic)
			return
		}
		t.Logf("expected error string: %v", err)
		return
	}

	topic = "not initialized"
	{
		v := V{}
		_, err := v.Insert(nil).After(0)
		shouldError(err)
		_, err = v.Insert(nil).Before(0)
		shouldError(err)
	}
	topic = "not array"
	{
		v := NewNull()
		_, err := v.Insert(nil).After(0)
		shouldError(err)
		_, err = v.Insert(nil).Before(0)
		shouldError(err)
	}
	topic = "param error"
	{
		v := NewArray()
		_, err := v.InsertNull().After(true)
		shouldError(err)
		_, err = v.InsertNull().Before(true)
		shouldError(err)
	}
	topic = "out of range"
	{
		v := NewArray()
		v.AppendNull().InTheEnd()
		v.AppendNull().InTheEnd()
		_, err := v.InsertNull().After(100)
		shouldError(err)
		_, err = v.InsertNull().Before(-100)
		shouldError(err)
	}
	topic = "deep not exist"
	{
		raw := `{"object":{"array":[1,2,3,4]}}`
		v, _ := UnmarshalString(raw)
		_, err := v.InsertNull().After("object", "not exist")
		shouldError(err)

		_, err = v.InsertNull().Before("object", "not exist")
		shouldError(err)
	}
}

func TestMiscAppendError(t *testing.T) {
	{
		v := V{}
		_, err := v.AppendString("blahblah").InTheBeginning()
		if err == nil {
			t.Errorf("expected error for an uninitialied object")
			return
		}

		_, err = v.AppendString("blahblah").InTheEnd()
		if err == nil {
			t.Errorf("expected error for an uninitialied object")
			return
		}
	}

	{
		v := NewString("blahblah")
		_, err := v.AppendString("blahblah").InTheBeginning()
		if err != ErrNotArrayValue {
			t.Errorf("expect error for ErrNotArrayValue")
			return
		}
	}

	{
		raw := `{"object":{"object":{"array":[[]],"object":{}}}}`
		v, err := UnmarshalString(raw)
		if err != nil {
			t.Errorf("UnmarshalString failed: %v", err)
			return
		}

		_, err = v.AppendNull().InTheBeginning("object", "object", "arrayNotExist")
		if err == nil {
			s, _ := v.MarshalString()
			t.Errorf("expect error for an inexist value, marshaled: %v", s)
			return
		}
		t.Logf("expected error: %v", err)

		_, err = v.AppendNull().InTheEnd("object", "object", "arrayNotExist")
		if err == nil {
			s, _ := v.MarshalString()
			t.Errorf("expect error for an inexist value, marshaled: %v", s)
			return
		}
		t.Logf("expected error: %v", err)

		_, err = v.AppendNull().InTheBeginning("object", "object")
		if err == nil {
			s, _ := v.MarshalString()
			t.Errorf("expect error for an not-arrayed value, marshaled: %v", s)
			return
		}
		t.Logf("expected error: %v", err)

		o, _ := v.Get("object", "object", "object")
		_, err = o.AppendNull().InTheEnd()
		if err == nil {
			s, _ := v.MarshalString()
			t.Errorf("expect error for an not-arrayed value, marshaled: %v", s)
			return
		}
		t.Logf("expected error: %v", err)

		_, err = v.InsertNull().Before("object", "object", "array", true)
		if err == nil {
			s, _ := v.MarshalString()
			t.Errorf("expect error for an invalid parameter, marshaled: %v", s)
			return
		}
		t.Logf("expected error: %v", err)
	}
}

func TestMiscDeleteError(t *testing.T) {
	var checkCount int
	shouldError := func(err error) {
		defer func() {
			checkCount++
		}()
		if err == nil {
			t.Errorf("%02d - error expected but not caught", checkCount)
			return
		}
		t.Logf("expected error string: %v", err)
		return
	}

	{
		var err error
		raw := `{"hello":"world","object":{"hello":"world","object":{"int":123456}},"array":[123456]}`
		v, _ := UnmarshalString(raw)

		// param error
		err = v.Delete(make(map[string]string))
		shouldError(err)

		// param error
		err = v.Delete("object", true)
		shouldError(err)

		// param error
		err = v.Delete("array", "2")
		shouldError(err)

		// not found error
		err = v.Delete("earth")
		shouldError(err)

		// out of range
		err = v.Delete("array", 100)
		shouldError(err)

		// not an object or array
		err = v.Delete("object", "object", "int", "number")
		shouldError(err)

		// not found error
		err = v.Delete("object", "bool", "string")
		shouldError(err)
	}
}
