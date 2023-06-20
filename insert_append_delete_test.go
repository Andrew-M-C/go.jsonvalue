package jsonvalue

import (
	"testing"
)

func testInsertAppendDelete(t *testing.T) {
	cv("insert/append", func() { testInsertAppend(t) })
	cv("must insert/append", func() { testMustInsertAppend(t) })
	cv("delete", func() { testDelete(t) })
	cv("must delete", func() { testMustDelete(t) })
	cv("test misc append functions", func() { testMiscAppend(t) })
	cv("test append and auto generate functions", func() { testAppendAndAutoGeneratePath(t) })
	cv("test misc insert functions", func() { testMiscInsert(t) })
	cv("test misc insert errors", func() { testMiscInsertError(t) })
	cv("test misc append errors", func() { testMiscAppendError(t) })
	cv("test misc delete errors", func() { testMiscDeleteError(t) })
}

func testInsertAppend(t *testing.T) {
	expected := `[123456,"hello","world",1234.123456789,true,["12345"],null,null,"MQ==",99999999]`
	a := NewArray()

	v, err := a.AppendString("world").InTheBeginning()
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `["world"]`)
	t.Log(a.MustMarshalString())

	v, err = a.AppendFloat64(1234.123456789).InTheEnd()
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `["world",1234.123456789]`)
	t.Log(a.MustMarshalString())

	v, err = a.InsertBool(true).After(-1)
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `["world",1234.123456789,true]`)
	t.Log(a.MustMarshalString())

	v, err = a.AppendNull().InTheEnd()
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `["world",1234.123456789,true,null]`)
	t.Log(a.MustMarshalString())

	v, err = a.InsertInt(123456).Before(0)
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `[123456,"world",1234.123456789,true,null]`)
	t.Log(a.MustMarshalString())

	v, err = a.InsertString("hello").After(0)
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,null]`)
	t.Log(a.MustMarshalString())

	v, err = a.InsertArray().After(-2)
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,[],null]`)
	t.Log(a.MustMarshalString())

	v, err = a.AppendString("12345").InTheEnd(-2)
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null]`)
	t.Log(a.MustMarshalString())

	v, err = a.Append(nil).InTheEnd()
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null,null]`)
	t.Log(a.MustMarshalString())

	v, err = a.AppendBytes([]byte("1")).InTheEnd()
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null,null,"MQ=="]`)
	t.Log(a.MustMarshalString())

	v, err = a.Append(99999999).InTheEnd()
	so(v, ne, nil)
	so(err, eq, nil)
	so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null,null,"MQ==",99999999]`)
	t.Log(a.MustMarshalString())

	s, _ := a.MarshalString()
	t.Logf("after SetXxx(): %v", s)

	so(s, eq, expected)

	// unmarshal and then marchal back
	a, err = UnmarshalString(expected)
	so(err, isNil)
	s, err = a.MarshalString()
	so(err, isNil)
	so(s, eq, expected)
}

func testMustInsertAppend(t *testing.T) {
	cv("general", func() {
		expected := `[123456,"hello","world",1234.123456789,true,["12345"],null,null,"MQ==",99999999]`
		a := NewArray()

		a.MustAppendString("world").InTheBeginning()
		so(a.MustMarshalString(), eq, `["world"]`)
		t.Log(a.MustMarshalString())

		a.MustAppendFloat64(1234.123456789).InTheEnd()
		so(a.MustMarshalString(), eq, `["world",1234.123456789]`)
		t.Log(a.MustMarshalString())

		a.MustInsertBool(true).After(-1)
		so(a.MustMarshalString(), eq, `["world",1234.123456789,true]`)
		t.Log(a.MustMarshalString())

		a.MustAppendNull().InTheEnd()
		so(a.MustMarshalString(), eq, `["world",1234.123456789,true,null]`)
		t.Log(a.MustMarshalString())

		a.MustInsertInt(123456).Before(0)
		so(a.MustMarshalString(), eq, `[123456,"world",1234.123456789,true,null]`)
		t.Log(a.MustMarshalString())

		a.MustInsertString("hello").After(0)
		so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,null]`)
		t.Log(a.MustMarshalString())

		a.MustInsertArray().After(-2)
		so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,[],null]`)
		t.Log(a.MustMarshalString())

		a.MustAppendString("12345").InTheEnd(-2)
		so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null]`)
		t.Log(a.MustMarshalString())

		a.MustAppend(nil).InTheEnd()
		so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null,null]`)
		t.Log(a.MustMarshalString())

		a.MustAppendBytes([]byte("1")).InTheEnd()
		so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null,null,"MQ=="]`)
		t.Log(a.MustMarshalString())

		a.MustAppend(99999999).InTheEnd()
		so(a.MustMarshalString(), eq, `[123456,"hello","world",1234.123456789,true,["12345"],null,null,"MQ==",99999999]`)
		t.Log(a.MustMarshalString())

		s, _ := a.MarshalString()
		t.Logf("after SetXxx(): %v", s)

		so(s, eq, expected)

		// unmarshal and then marchal back
		a, err := UnmarshalString(expected)
		so(err, isNil)
		s, err = a.MarshalString()
		so(err, isNil)
		so(s, eq, expected)
	})

	cv("MustInsertInt64", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertInt64(10).Before(1)
		b.MustInsertInt64(10).Before(1)
		so(a.MustMarshalString(), eq, `[1,10,2]`)
		so(b.MustMarshalString(), eq, `[1,10,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertInt32", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertInt32(-1).Before(1)
		b.MustInsertInt32(-1).Before(1)
		so(a.MustMarshalString(), eq, `[1,-1,2]`)
		so(b.MustMarshalString(), eq, `[1,-1,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertUint", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertUint(10).Before(1)
		b.MustInsertUint(10).Before(1)
		so(a.MustMarshalString(), eq, `[1,10,2]`)
		so(b.MustMarshalString(), eq, `[1,10,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertUint64", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertUint64(10).Before(1)
		b.MustInsertUint64(10).Before(1)
		so(a.MustMarshalString(), eq, `[1,10,2]`)
		so(b.MustMarshalString(), eq, `[1,10,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertUint32", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertUint32(10).Before(1)
		b.MustInsertUint32(10).Before(1)
		so(a.MustMarshalString(), eq, `[1,10,2]`)
		so(b.MustMarshalString(), eq, `[1,10,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertFloat64", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertFloat64(-1.5).Before(1)
		b.MustInsertFloat64(-1.5).Before(1)
		so(a.MustMarshalString(), eq, `[1,-1.5,2]`)
		so(b.MustMarshalString(), eq, `[1,-1.5,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertFloat32", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertFloat32(-1.5).Before(1)
		b.MustInsertFloat32(-1.5).Before(1)
		so(a.MustMarshalString(), eq, `[1,-1.5,2]`)
		so(b.MustMarshalString(), eq, `[1,-1.5,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertNull", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertNull().Before(1)
		b.MustInsertNull().Before(1)
		so(a.MustMarshalString(), eq, `[1,null,2]`)
		so(b.MustMarshalString(), eq, `[1,null,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustInsertObject", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.InsertObject().After(1)
		b.MustInsertObject().After(1)
		so(a.MustMarshalString(), eq, `[1,2,{}]`)
		so(b.MustMarshalString(), eq, `[1,2,{}]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendBool", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendBool(true).InTheBeginning()
		b.MustAppendBool(true).InTheBeginning()
		so(a.MustMarshalString(), eq, `[true,1,2]`)
		so(b.MustMarshalString(), eq, `[true,1,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendInt64", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendInt64(-100).InTheEnd()
		b.MustAppendInt64(-100).InTheEnd()
		so(a.MustMarshalString(), eq, `[1,2,-100]`)
		so(b.MustMarshalString(), eq, `[1,2,-100]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendInt32", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendInt32(-100).InTheBeginning()
		b.MustAppendInt32(-100).InTheBeginning()
		so(a.MustMarshalString(), eq, `[-100,1,2]`)
		so(b.MustMarshalString(), eq, `[-100,1,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendUint", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendUint(100).InTheBeginning()
		b.MustAppendUint(100).InTheBeginning()
		so(a.MustMarshalString(), eq, `[100,1,2]`)
		so(b.MustMarshalString(), eq, `[100,1,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendUint64", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendUint64(100).InTheBeginning()
		b.MustAppendUint64(100).InTheBeginning()
		so(a.MustMarshalString(), eq, `[100,1,2]`)
		so(b.MustMarshalString(), eq, `[100,1,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendUint32", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendUint32(100).InTheBeginning()
		b.MustAppendUint32(100).InTheBeginning()
		so(a.MustMarshalString(), eq, `[100,1,2]`)
		so(b.MustMarshalString(), eq, `[100,1,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendFloat64", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendFloat64(1.25).InTheBeginning()
		b.MustAppendFloat64(1.25).InTheBeginning()
		so(a.MustMarshalString(), eq, `[1.25,1,2]`)
		so(b.MustMarshalString(), eq, `[1.25,1,2]`)
		so(a.Equal(b), isTrue)
	})

	cv("MustAppendFloat32", func() {
		a := New([]int{1, 2})
		b := New([]int{1, 2})
		_, _ = a.AppendFloat32(-1.25).InTheBeginning()
		b.MustAppendFloat32(-1.25).InTheBeginning()
		so(a.MustMarshalString(), eq, `[-1.25,1,2]`)
		so(b.MustMarshalString(), eq, `[-1.25,1,2]`)
		so(a.Equal(b), isTrue)
	})
}

func testDelete(t *testing.T) {
	raw := `{"array":[1,2,3,4,5,6],"string":"string to be deleted","object":{"number":12345},"OBJECT":{}}`
	o, err := UnmarshalString(raw)
	so(err, isNil)

	s, _ := o.MarshalString()
	t.Logf("parsed object: %v", s)

	err = o.Delete("oBJECT") // this key not exists
	so(err, isErr)

	err = o.Delete("object", "number")
	so(err, isNil)

	sub, err := o.Get("object")
	so(err, isNil)

	s, _ = sub.MarshalString()
	so(s, eq, "{}")

	err = o.Delete("object", "number")
	so(err, isErr, ErrNotFound)

	err = o.Delete("object")
	so(err, isNil)

	_, err = o.Caseless().Get("object")
	so(err, isNil)

	err = o.Delete("object")
	so(err, isErr)

	err = o.Caseless().Delete("object") // delete another "object", actually "OBJECT"
	so(err, isNil)

	err = o.Caseless().Delete("object") // delete again
	so(err, isErr)

	err = o.Caseless().Delete("NOT_EXIST")
	so(err, isErr)

	_, err = o.Get("object")
	so(err, isErr, ErrNotFound)

	err = o.Delete("string")
	so(err, isNil)

	s, _ = o.MarshalString()
	so(s, eq, `{"array":[1,2,3,4,5,6]}`)

	err = o.Delete("array", 1)
	so(err, isNil)

	s, _ = o.MarshalString()
	so(s, eq, `{"array":[1,3,4,5,6]}`)
}

func testMustDelete(t *testing.T) {
	raw := `{"array":[1,2,3,4,5,6],"string":"string to be deleted","object":{"number":12345},"OBJECT":{}}`
	o, err := UnmarshalString(raw)
	so(err, isNil)

	s, _ := o.MarshalString()
	t.Logf("parsed object: %v", s)

	o.MustDelete("oBJECT") // this key not exists
	err = o.Delete("object", "number")
	so(err, isNil)

	sub, err := o.Get("object")
	so(err, isNil)

	s, _ = sub.MarshalString()
	so(s, eq, "{}")

	o.MustDelete("object")
	_, err = o.Caseless().Get("object")
	so(err, isNil)

	o.MustDelete("object")
	o.Caseless().MustDelete("object")   // delete another "object", actually "OBJECT"
	err = o.Caseless().Delete("object") // delete again
	so(err, isErr)

	err = o.Caseless().Delete("NOT_EXIST")
	so(err, isErr)

	_, err = o.Get("object")
	so(err, isErr, ErrNotFound)

	o.MustDelete("string")
	s, _ = o.MarshalString()
	so(s, eq, `{"array":[1,2,3,4,5,6]}`)

	o.MustDelete("array", 1)
	s, _ = o.MarshalString()
	so(s, eq, `{"array":[1,3,4,5,6]}`)
}

func testMiscAppend(t *testing.T) {
	expected := `[true,-1,2,-3,4,-5,6,-7.7,8.8000,{},[[false],null]]`
	a := NewArray()
	_, _ = a.AppendBool(true).InTheBeginning()
	_, _ = a.AppendInt(-1).InTheEnd()
	_, _ = a.AppendUint(2).InTheEnd()
	_, _ = a.AppendInt32(-3).InTheEnd()
	_, _ = a.AppendUint32(4).InTheEnd()
	_, _ = a.AppendInt64(-5).InTheEnd()
	_, _ = a.AppendUint64(6).InTheEnd()
	_, _ = a.AppendFloat32(-7.7).InTheEnd()
	_, _ = a.Append(NewFloat64f(8.8, 'f', 4)).InTheEnd()
	_, _ = a.AppendObject().InTheEnd()
	_, _ = a.AppendArray().InTheEnd()
	_, _ = a.AppendNull().InTheEnd(-1)
	_, _ = a.AppendArray().InTheBeginning(-1)
	_, _ = a.AppendBool(false).InTheBeginning(-1, 0)

	s, _ := a.MarshalString()
	so(s, eq, expected)
}

func testAppendAndAutoGeneratePath(t *testing.T) {
	expected := `{"arr":[1]}`

	o := NewObject()
	_, err := o.AppendInt(1).InTheEnd("arr")
	so(err, isNil)

	so(o.MustMarshalString(), eq, expected)
}

func testMiscInsert(t *testing.T) {
	expected := `[null,1,-2,3,-4,5,-6,7.7,-8.88888,true,false,null,null,{},-2,"insert test",[[null,-11,22]]]`

	var err error
	var c *V
	v := NewArray()

	_, err = v.InsertNull().Before(0)
	so(err, isErr)

	_, err = v.InsertNull().After(0)
	so(err, isErr)

	_, err = v.AppendNull().InTheBeginning()
	so(err, isNil)

	_, err = v.InsertNull().Before(10000)
	so(err, isErr)

	_, err = v.InsertNull().After(10000)
	so(err, isErr)

	_, err = v.InsertNull().Before(-10000)
	so(err, isErr)

	_, err = v.InsertNull().After(-10000)
	so(err, isErr)

	_, err = v.InsertNull().Before(-2)
	so(err, isErr)

	c, err = v.InsertUint(1).After(-1)
	so(err, isNil)
	so(c.Int(), eq, 1)

	c, err = v.InsertInt(-2).After(-1)
	so(err, isNil)
	so(c.Int(), eq, -2)

	c, err = v.InsertUint64(3).After(-1)
	so(err, isNil)
	so(c.Int(), eq, 3)

	c, err = v.InsertInt64(-4).After(-1)
	so(err, isNil)
	so(c.Int(), eq, -4)

	c, err = v.InsertUint32(5).After(-1)
	so(err, isNil)
	so(c.Int(), eq, 5)

	c, err = v.InsertInt32(-6).After(-1)
	so(err, isNil)
	so(c.Int(), eq, -6)

	c, err = v.InsertFloat32(7.7).After(-1)
	t.Log(v.MustMarshalString())
	so(err, isNil)
	so(c.String(), eq, "7.7")

	c, err = v.InsertFloat64(-8.88888).After(-1)
	so(err, isNil)
	so(c.Float64(), eq, -8.88888)

	c, err = v.InsertBool(true).After(-1)
	so(err, isNil)
	so(c.Bool(), isTrue)

	c, err = v.InsertBool(false).After(-1)
	so(err, isNil)
	so(c.IsBoolean(), isTrue)
	so(c.Bool(), isFalse)

	c, err = v.Insert(nil).After(-1)
	so(err, isNil)
	so(c.IsNull(), isTrue)

	c, err = v.InsertNull().After(-1)
	so(err, isNil)
	so(c.IsNull(), isTrue)

	c, err = v.InsertObject().After(-1)
	so(err, isNil)
	so(c.IsObject(), isTrue)

	c, err = v.InsertArray().After(-1)
	so(err, isNil)
	so(c.IsArray(), isTrue)

	c, err = v.AppendArray().InTheBeginning(-1)
	so(err, isNil)
	so(c.IsArray(), isTrue)

	c, err = v.AppendInt(-11).InTheBeginning(-1, 0)
	so(err, isNil)
	so(c.Int(), eq, -11)

	c, err = v.InsertUint(22).After(-1, 0, 0)
	so(err, isNil)
	so(c.Int(), eq, 22)

	c, err = v.InsertNull().Before(-1, 0, 0)
	so(err, isNil)
	so(c.IsNull(), isTrue)

	c, err = v.InsertInt(-2).Before(-1)
	so(err, isNil)
	so(c.Int(), eq, -2)

	c, err = v.Insert("insert test").Before(-1)
	so(err, isNil)
	so(c.String(), eq, "insert test")

	s, _ := v.MarshalString()
	so(s, eq, expected)
}

func testMiscInsertError(t *testing.T) {
	cv("not initialized", func() {
		v := V{}
		_, err := v.Insert(nil).After(0)
		so(err, isErr)
		_, err = v.Insert(nil).Before(0)
		so(err, isErr)
	})

	cv("not array", func() {
		v := NewNull()
		_, err := v.Insert(nil).After(0)
		so(err, isErr)
		_, err = v.Insert(nil).Before(0)
		so(err, isErr)
	})

	cv("param error", func() {
		v := NewArray()
		_, err := v.InsertNull().After(true)
		so(err, isErr)
		_, err = v.InsertNull().Before(true)
		so(err, isErr)
	})

	cv("out of range", func() {
		v := NewArray()
		v.MustAppendNull().InTheEnd()
		v.MustAppendNull().InTheEnd()
		_, err := v.InsertNull().After(100)
		so(err, isErr)
		_, err = v.InsertNull().Before(-100)
		so(err, isErr)
	})

	cv("deep not exist", func() {
		raw := `{"object":{"array":[1,2,3,4]}}`
		v, _ := UnmarshalString(raw)
		_, err := v.InsertNull().After("object", "not exist")
		so(err, isErr)

		_, err = v.InsertNull().Before("object", "not exist")
		so(err, isErr)
	})

	cv("invalid insert type", func() {
		a := MustUnmarshalString(`[0]`)
		ch := make(chan struct{}, 1)
		_, err := a.Insert(ch).After(0)
		so(err, isErr)

		_, err = a.Insert(ch).Before(0)
		so(err, isErr)
	})
}

func testMiscAppendError(t *testing.T) {
	cv("uninitialized AppendString to uninitialized V", func() {
		v := V{}
		_, err := v.AppendString("blahblah").InTheBeginning()
		so(err, isErr)

		_, err = v.AppendString("blahblah").InTheEnd()
		so(err, isErr)
	})

	cv("uninitialized AppendString to string", func() {
		v := NewString("blahblah")
		_, err := v.AppendString("blahblah").InTheBeginning()
		so(err, isErr)
	})

	cv("append non exist data", func() {
		raw := `{"object":{"object":{"array":[[]],"object":{}}}}`
		v, err := UnmarshalString(raw)
		so(err, isNil)

		_, err = v.AppendNull().InTheBeginning("object", "arrayNotExist", "arrayNotExistForTheBeginning")
		so(err, isNil)

		_, err = v.AppendNull().InTheEnd("object", "arrayNotExist", "arrayNotExistForTheEnd")
		so(err, isNil)

		_, err = v.AppendNull().InTheBeginning("object", "object")
		so(err, isErr)

		_, err = v.AppendNull().InTheEnd("object", "object")
		so(err, isErr)

		err = v.GetNull("object", "arrayNotExist", "arrayNotExistForTheBeginning", 0)
		so(err, isNil)

		err = v.GetNull("object", "arrayNotExist", "arrayNotExistForTheEnd", 0)
		so(err, isNil)
	})

	cv("append/insert to error type", func() {
		raw := `{"object":{"object":{"array":[[]],"object":{}}}}`
		v, err := UnmarshalString(raw)
		so(err, isNil)

		_, err = v.AppendNull().InTheBeginning("object", "object")
		so(err, isErr)

		_, err = v.AppendNull().InTheBeginning("object", "object")
		so(err, isErr)

		_, err = v.InsertNull().After("object", "object", "object", 0)
		so(err, isErr)
		t.Logf("expected error: %v", err)
	})

	cv("insert non-exist data", func() {
		raw := `{"object":{"object":{"array":[[]],"object":{}}}}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		o, _ := v.Get("object", "object", "object")
		_, err = o.AppendNull().InTheEnd()
		so(err, isErr)
		t.Logf("expected error: %v", err)

		_, err = v.InsertNull().Before("object", "object", "array", true)
		so(err, isErr)
		t.Logf("expected error: %v", err)
	})
}

func testMiscDeleteError(t *testing.T) {
	var err error
	raw := `{"Hello":"world","object":{"hello":"world","object":{"int":123456}},"array":[123456]}`
	v, _ := UnmarshalString(raw)

	// param error
	err = v.Delete(make(map[string]string))
	so(err, isErr)

	// param error
	err = v.Delete("object", true)
	so(err, isErr)

	// param error
	err = v.Delete("array", "2")
	so(err, isErr)

	// not found error
	err = v.Delete("earth")
	so(err, isErr)

	// out of range
	err = v.Delete("array", 100)
	so(err, isErr)

	// not an object or array
	err = v.Delete("object", "object", "int", "number")
	so(err, isErr)

	// not found error
	err = v.Delete("object", "bool", "string")
	so(err, isErr)
}
