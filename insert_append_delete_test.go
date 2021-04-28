package jsonvalue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestInsertAppendDelete(t *testing.T) {
	test(t, "insert/append", testInsertAppend)
	test(t, "delete", testDelete)
	test(t, "test misc append functions", testMiscAppend)
	test(t, "test misc insert functions", testMiscInsert)
	test(t, "test misc insert errors", testMiscInsertError)
	test(t, "test misc append errors", testMiscAppendError)
	test(t, "test misc delete errors", testMiscDeleteError)
}

func testInsertAppend(t *testing.T) {
	expected := `[123456,"hello","world",1234.123456789,true,["12345"],null,null]`
	a := NewArray()

	a.AppendString("world").InTheBeginning()
	So(a.MustMarshalString(), ShouldEqual, `["world"]`)
	t.Log(a.MustMarshalString())

	a.AppendFloat64(1234.123456789, 9).InTheEnd()
	So(a.MustMarshalString(), ShouldEqual, `["world",1234.123456789]`)
	t.Log(a.MustMarshalString())

	a.InsertBool(true).After(-1)
	So(a.MustMarshalString(), ShouldEqual, `["world",1234.123456789,true]`)
	t.Log(a.MustMarshalString())

	a.AppendNull().InTheEnd()
	So(a.MustMarshalString(), ShouldEqual, `["world",1234.123456789,true,null]`)
	t.Log(a.MustMarshalString())

	a.InsertInt(123456).Before(0)
	So(a.MustMarshalString(), ShouldEqual, `[123456,"world",1234.123456789,true,null]`)
	t.Log(a.MustMarshalString())

	a.InsertString("hello").After(0)
	So(a.MustMarshalString(), ShouldEqual, `[123456,"hello","world",1234.123456789,true,null]`)
	t.Log(a.MustMarshalString())

	a.InsertArray().After(-2)
	So(a.MustMarshalString(), ShouldEqual, `[123456,"hello","world",1234.123456789,true,[],null]`)
	t.Log(a.MustMarshalString())

	a.AppendString("12345").InTheEnd(-2)
	So(a.MustMarshalString(), ShouldEqual, `[123456,"hello","world",1234.123456789,true,["12345"],null]`)
	t.Log(a.MustMarshalString())

	a.Append(nil).InTheEnd()
	So(a.MustMarshalString(), ShouldEqual, `[123456,"hello","world",1234.123456789,true,["12345"],null,null]`)
	t.Log(a.MustMarshalString())

	s, _ := a.MarshalString()
	t.Logf("after SetXxx(): %v", s)

	So(s, ShouldEqual, expected)

	// unmarshal and then marchal back
	a, err := UnmarshalString(expected)
	So(err, ShouldBeNil)
	s, err = a.MarshalString()
	So(err, ShouldBeNil)
	So(s, ShouldEqual, expected)
}

func testDelete(t *testing.T) {
	raw := `{"array":[1,2,3,4,5,6],"string":"string to be deleted","object":{"number":12345},"OBJECT":{}}`
	o, err := UnmarshalString(raw)
	So(err, ShouldBeNil)

	s, _ := o.MarshalString()
	t.Logf("parsed object: %v", s)

	err = o.Delete("oBJECT") // this key not exists
	So(err, ShouldBeError)

	err = o.Delete("object", "number")
	So(err, ShouldBeNil)

	sub, err := o.Get("object")
	So(err, ShouldBeNil)

	s, _ = sub.MarshalString()
	So(s, ShouldEqual, "{}")

	err = o.Delete("object", "number")
	So(err, ShouldBeError, ErrNotFound)

	err = o.Delete("object")
	So(err, ShouldBeNil)

	_, err = o.Caseless().Get("object")
	So(err, ShouldBeNil)

	err = o.Delete("object")
	So(err, ShouldBeError)

	err = o.Caseless().Delete("object") // delete another "object", actually "OBJECT"
	So(err, ShouldBeNil)

	err = o.Caseless().Delete("object") // delete again
	So(err, ShouldBeError)

	err = o.Caseless().Delete("NOT_EXIST")
	So(err, ShouldBeError)

	_, err = o.Get("object")
	So(err, ShouldBeError, ErrNotFound)

	err = o.Delete("string")
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	So(s, ShouldEqual, `{"array":[1,2,3,4,5,6]}`)

	err = o.Delete("array", 1)
	So(err, ShouldBeNil)

	s, _ = o.MarshalString()
	So(s, ShouldEqual, `{"array":[1,3,4,5,6]}`)
}

func testMiscAppend(t *testing.T) {
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
	So(s, ShouldEqual, expected)
}

func testMiscInsert(t *testing.T) {
	expected := `[null,1,-2,3,-4,5,-6,7.7,-8.88888,true,false,null,null,{},-2,[[null,-11,22]]]`

	var err error
	var c *V
	v := NewArray()

	_, err = v.InsertNull().Before(0)
	So(err, ShouldBeError)

	_, err = v.InsertNull().After(0)
	So(err, ShouldBeError)

	_, err = v.AppendNull().InTheBeginning()
	So(err, ShouldBeNil)

	_, err = v.InsertNull().Before(10000)
	So(err, ShouldBeError)

	_, err = v.InsertNull().After(10000)
	So(err, ShouldBeError)

	_, err = v.InsertNull().Before(-10000)
	So(err, ShouldBeError)

	_, err = v.InsertNull().After(-10000)
	So(err, ShouldBeError)

	_, err = v.InsertNull().Before(-2)
	So(err, ShouldBeError)

	c, err = v.InsertUint(1).After(-1)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, 1)

	c, err = v.InsertInt(-2).After(-1)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, -2)

	c, err = v.InsertUint64(3).After(-1)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, 3)

	c, err = v.InsertInt64(-4).After(-1)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, -4)

	c, err = v.InsertUint32(5).After(-1)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, 5)

	c, err = v.InsertInt32(-6).After(-1)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, -6)

	c, err = v.InsertFloat32(7.7, 1).After(-1)
	Print(v.MustMarshalString())
	So(err, ShouldBeNil)
	So(c.String(), ShouldEqual, "7.7")

	c, err = v.InsertFloat64(-8.88888, 5).After(-1)
	So(err, ShouldBeNil)
	So(c.Float64(), ShouldEqual, -8.88888)

	c, err = v.InsertBool(true).After(-1)
	So(err, ShouldBeNil)
	So(c.Bool(), ShouldBeTrue)

	c, err = v.InsertBool(false).After(-1)
	So(err, ShouldBeNil)
	So(c.IsBoolean(), ShouldBeTrue)
	So(c.Bool(), ShouldBeFalse)

	c, err = v.Insert(nil).After(-1)
	So(err, ShouldBeNil)
	So(c.IsNull(), ShouldBeTrue)

	c, err = v.InsertNull().After(-1)
	So(err, ShouldBeNil)
	So(c.IsNull(), ShouldBeTrue)

	c, err = v.InsertObject().After(-1)
	So(err, ShouldBeNil)
	So(c.IsObject(), ShouldBeTrue)

	c, err = v.InsertArray().After(-1)
	So(err, ShouldBeNil)
	So(c.IsArray(), ShouldBeTrue)

	c, err = v.AppendArray().InTheBeginning(-1)
	So(err, ShouldBeNil)
	So(c.IsArray(), ShouldBeTrue)

	c, err = v.AppendInt(-11).InTheBeginning(-1, 0)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, -11)

	c, err = v.InsertUint(22).After(-1, 0, 0)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, 22)

	c, err = v.InsertNull().Before(-1, 0, 0)
	So(err, ShouldBeNil)
	So(c.IsNull(), ShouldBeTrue)

	c, err = v.InsertInt(-2).Before(-1)
	So(err, ShouldBeNil)
	So(c.Int(), ShouldEqual, -2)

	s, _ := v.MarshalString()
	So(s, ShouldEqual, expected)
}

func testMiscInsertError(t *testing.T) {
	Convey("not initialized", func() {
		v := V{}
		_, err := v.Insert(nil).After(0)
		So(err, ShouldBeError)
		_, err = v.Insert(nil).Before(0)
		So(err, ShouldBeError)
	})

	Convey("not array", func() {
		v := NewNull()
		_, err := v.Insert(nil).After(0)
		So(err, ShouldBeError)
		_, err = v.Insert(nil).Before(0)
		So(err, ShouldBeError)
	})

	Convey("param error", func() {
		v := NewArray()
		_, err := v.InsertNull().After(true)
		So(err, ShouldBeError)
		_, err = v.InsertNull().Before(true)
		So(err, ShouldBeError)
	})

	Convey("out of range", func() {
		v := NewArray()
		v.AppendNull().InTheEnd()
		v.AppendNull().InTheEnd()
		_, err := v.InsertNull().After(100)
		So(err, ShouldBeError)
		_, err = v.InsertNull().Before(-100)
		So(err, ShouldBeError)
	})

	Convey("deep not exist", func() {
		raw := `{"object":{"array":[1,2,3,4]}}`
		v, _ := UnmarshalString(raw)
		_, err := v.InsertNull().After("object", "not exist")
		So(err, ShouldBeError)

		_, err = v.InsertNull().Before("object", "not exist")
		So(err, ShouldBeError)
	})

	Convey("uninitialized append", func() {
		_, err := (&Append{}).InTheBeginning("dummy")
		So(err, ShouldBeError)
	})
}

func testMiscAppendError(t *testing.T) {
	Convey("uninitialized AppendString to uninitialized V", func() {
		v := V{}
		_, err := v.AppendString("blahblah").InTheBeginning()
		So(err, ShouldBeError)

		_, err = v.AppendString("blahblah").InTheEnd()
		So(err, ShouldBeError)
	})

	Convey("uninitialized AppendString to string", func() {
		v := NewString("blahblah")
		_, err := v.AppendString("blahblah").InTheBeginning()
		So(err, ShouldBeError)
	})

	Convey("misc error", func() {
		raw := `{"object":{"object":{"array":[[]],"object":{}}}}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)

		_, err = v.AppendNull().InTheBeginning("object", "object", "arrayNotExist")
		So(err, ShouldBeError)
		t.Logf("expected error: %v", err)

		_, err = v.AppendNull().InTheEnd("object", "object", "arrayNotExist")
		So(err, ShouldBeError)
		t.Logf("expected error: %v", err)

		_, err = v.AppendNull().InTheBeginning("object", "object")
		So(err, ShouldBeError)
		t.Logf("expected error: %v", err)

		o, _ := v.Get("object", "object", "object")
		_, err = o.AppendNull().InTheEnd()
		So(err, ShouldBeError)
		t.Logf("expected error: %v", err)

		_, err = v.InsertNull().Before("object", "object", "array", true)
		So(err, ShouldBeError)
		t.Logf("expected error: %v", err)
	})
}

func testMiscDeleteError(t *testing.T) {
	var err error
	raw := `{"Hello":"world","object":{"hello":"world","object":{"int":123456}},"array":[123456]}`
	v, _ := UnmarshalString(raw)

	// param error
	err = v.Delete(make(map[string]string))
	So(err, ShouldBeError)

	// param error
	err = v.Delete("object", true)
	So(err, ShouldBeError)

	// param error
	err = v.Delete("array", "2")
	So(err, ShouldBeError)

	// not found error
	err = v.Delete("earth")
	So(err, ShouldBeError)

	// out of range
	err = v.Delete("array", 100)
	So(err, ShouldBeError)

	// not an object or array
	err = v.Delete("object", "object", "int", "number")
	So(err, ShouldBeError)

	// not found error
	err = v.Delete("object", "bool", "string")
	So(err, ShouldBeError)
}
