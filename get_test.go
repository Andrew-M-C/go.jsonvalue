package jsonvalue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGet(t *testing.T) {
	test(t, "basic Get function", testGet)
	test(t, "misc get errors", testMiscError)
	test(t, "caseless get", testCaselessGet)
}

func testGet(t *testing.T) {
	full := `{"data":{"message":["hello","world",true,null],"author":"Andrew","year":2019,"YYYY.MM":2019.12,"negative":-1234,"num_in_str":"2020.02","negative_in_str":"-12345","invalid_num_in_str":"2020/02"}}`

	o, err := UnmarshalString(full)
	So(err, ShouldBeNil)

	b, _ := o.Marshal()
	t.Logf("unmarshal back: %s", string(b))

	Convey("general GetString", func() {
		s, err := o.GetString("data", "author")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "Andrew")
	})

	Convey("GetInt", func() {
		i, err := o.GetInt("data", "year")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, 2019)
	})

	Convey("GetUint", func() {
		i, err := o.GetUint("data", "year")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, 2019)
	})

	Convey("GetInt64", func() {
		i, err := o.GetInt64("data", "year")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, 2019)
	})

	Convey("GetUint64", func() {
		i, err := o.GetUint64("data", "year")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, 2019)
	})

	Convey("GetInt32", func() {
		i, err := o.GetInt32("data", "negative")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, -1234)
	})

	Convey("Caseless.GetInt32", func() {
		i, err := o.Caseless().GetInt32("data", "negATive")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, -1234)
	})

	Convey("GetInt32 but not caseless", func() {
		_, err := o.GetInt32("data", "negATive")
		So(err, ShouldBeError)
	})

	Convey("GetUint32", func() {
		i, err := o.GetUint32("data", "year")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, 2019)
	})

	Convey("GetFloat64", func() {
		f, err := o.GetFloat64("data", "YYYY.MM")
		So(err, ShouldBeNil)
		So(f, ShouldEqual, 2019.12)
	})

	Convey("GetFloat32", func() {
		f, err := o.GetFloat32("data", "YYYY.MM")
		So(err, ShouldBeNil)
		So(f, ShouldEqual, 2019.12)
	})

	Convey("GetNull", func() {
		err := o.GetNull("data", "message", -1)
		So(err, ShouldBeNil)
	})

	Convey("GetBool", func() {
		b, err := o.GetBool("data", "message", 2)
		So(err, ShouldBeNil)
		So(b, ShouldBeTrue)
	})

	Convey("GetString from array of first one", func() {
		s, err := o.GetString("data", "message", 0)
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "hello")
	})

	Convey("GetString from array of last third one", func() {
		s, err := o.GetString("data", "message", -3)
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "world")
	})

	Convey("Len", func() {
		l := o.Len()
		So(l, ShouldEqual, 1)

		v, _ := o.Get("data", "message")
		l = v.Len()
		So(l, ShouldEqual, 4)

		v, _ = o.Get("data", "author")
		l = v.Len()
		So(l, ShouldEqual, 0)
	})

	Convey("GetObject", func() {
		v, err := o.GetObject("data")
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)
	})

	Convey("GetObject in object", func() {
		v, err := o.Caseless().GetObject("Data")
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)
	})

	Convey("nil V string", func() {
		v, _ := o.GetObject("not_exist")
		So(v.String(), ShouldEqual, "")
	})

	Convey("key: num_in_str", func() {
		v, err := o.Get("data", "num_in_str")
		So(err, ShouldBeNil)
		So(v.Int(), ShouldEqual, 2020)
		So(v.Float64(), ShouldEqual, 2020.02)
		So(v.String(), ShouldEqual, "2020.02")
	})

	Convey("key: negative_in_str", func() {
		v, err := o.Get("data", "negative_in_str")
		So(err, ShouldBeNil)
		So(v.Int(), ShouldEqual, -12345)
		So(v.IsString(), ShouldBeTrue)
		So(v.String(), ShouldEqual, "-12345")
	})

	Convey("key: invalid_num_in_str", func() {
		v, err := o.Get("data", "invalid_num_in_str")
		So(err, ShouldBeNil)
		So(v.IsString(), ShouldBeTrue)
		So(v.Float64(), ShouldEqual, 0)
	})
}

func testMiscError(t *testing.T) {
	var err error
	raw := `{"array":[0,1,2,3],"string":"hello, world","number":1234.12345}`
	v, _ := UnmarshalString(raw)

	// param error
	_, err = v.GetInt("array", true)
	So(err, ShouldNotBeNil)
	_, err = v.GetString(true)
	So(err, ShouldNotBeNil)

	// Caseless via non object or array
	child, err := v.Get("string")
	So(err, ShouldBeNil)
	_, err = child.Caseless().Get("NOT_EXIST")
	So(err, ShouldBeError)

	// out of range
	_, err = v.Get("array", 100)
	So(err, ShouldNotBeNil)

	// not support
	_, err = v.Get("string", "hello")
	So(err, ShouldNotBeNil)

	// GetString
	_, err = v.GetString("number")
	So(err, ShouldNotBeNil)
	_, err = v.GetString("not exist")
	So(err, ShouldNotBeNil)

	// GetInt... and GetUint..
	_, err = v.GetInt("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetUint("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetInt64("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetUint64("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetInt32("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetUint32("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetFloat64("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetFloat32("string")
	So(err, ShouldNotBeNil)

	// number not exist
	_, err = v.GetString("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetInt("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetUint("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetInt64("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetUint64("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetInt32("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetUint32("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetFloat64("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetFloat32("not exist")
	So(err, ShouldNotBeNil)

	// GetObject and GetArray
	_, err = v.GetObject("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetArray("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetObject("not exist")
	So(err, ShouldNotBeNil)
	_, err = v.GetArray("not exist")
	So(err, ShouldNotBeNil)

	// GetBool and GetNull
	_, err = v.GetBool("string")
	So(err, ShouldNotBeNil)
	err = v.GetNull("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetBool("not exist")
	So(err, ShouldNotBeNil)
	err = v.GetNull("not exist")
	So(err, ShouldNotBeNil)

	// GetBytes
	_, err = v.GetBytes("string")
	So(err, ShouldNotBeNil)
	_, err = v.GetBytes("array")
	So(err, ShouldNotBeNil)
}

func testCaselessGet(t *testing.T) {
	raw := `{"data":{"STRING":"hello, world","INTEGER":12345,"TRUE":true,"FALSE":false,"NULL":null,"FLOAT":1234.5678,"OBJECT":{},"ARRAY":[]}}`

	v, err := UnmarshalString(raw)
	v.SetBytes([]byte{1, 2, 3, 4}).At("data", "BYTES")

	t.Log(v.MustMarshalString())

	So(err, ShouldBeNil)
	So(v.IsObject(), ShouldBeTrue)

	_, err = v.Get("data", "object")
	So(err, ShouldBeError)
	_, err = v.Caseless().Get("data", "object")
	So(err, ShouldBeNil)

	_, err = v.GetObject("data", "object")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetObject("data", "object")
	So(err, ShouldBeNil)

	_, err = v.GetArray("data", "array")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetArray("data", "array")
	So(err, ShouldBeNil)

	_, err = v.GetBytes("data", "bytes")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetBytes("data", "bytes")
	So(err, ShouldBeNil)

	_, err = v.GetString("data", "string")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetString("data", "string")
	So(err, ShouldBeNil)

	_, err = v.GetInt("data", "integer")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetInt("data", "integer")
	So(err, ShouldBeNil)

	_, err = v.GetUint("data", "integer")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetUint("data", "integer")
	So(err, ShouldBeNil)

	_, err = v.GetInt64("data", "integer")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetInt64("data", "integer")
	So(err, ShouldBeNil)

	_, err = v.GetUint64("data", "integer")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetUint64("data", "integer")
	So(err, ShouldBeNil)

	_, err = v.GetInt32("data", "integer")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetInt32("data", "integer")
	So(err, ShouldBeNil)

	_, err = v.GetUint32("data", "integer")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetUint32("data", "integer")
	So(err, ShouldBeNil)

	_, err = v.GetFloat64("data", "float")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetFloat64("data", "float")
	So(err, ShouldBeNil)

	_, err = v.GetFloat32("data", "float")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetFloat32("data", "float")
	So(err, ShouldBeNil)

	_, err = v.GetBool("data", "true")
	So(err, ShouldBeError)
	_, err = v.Caseless().GetBool("data", "true")
	So(err, ShouldBeNil)

	err = v.GetNull("data", "null")
	So(err, ShouldBeError)
	err = v.Caseless().GetNull("data", "null")
	So(err, ShouldBeNil)
}
