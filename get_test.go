package jsonvalue

import (
	"errors"
	"testing"
)

func testGet(t *testing.T) {
	cv("basic Get function", func() { testJsonvalue_Get(t) })
	cv("misc get errors", func() { testMiscError(t) })
	cv("caseless get", func() { testCaselessGet(t) })
	cv("NotExist get", func() { testNotExistGet(t) })
	cv("get number from a string", func() { testGetNumFromString(t) })
}

func testJsonvalue_Get(t *testing.T) {
	full := `{"data":{"message":["hello","world",true,null],"author":"Andrew","year":2019,"YYYY.MM":2019.12,"negative":-1234,"num_in_str":"2020.02","negative_in_str":"-12345","invalid_num_in_str":"2020/02"}}`

	o, err := UnmarshalString(full)
	so(err, isNil)

	b, _ := o.Marshal()
	t.Logf("unmarshal back: %s", string(b))

	cv("general GetString", func() {
		s, err := o.GetString("data", "author")
		so(err, isNil)
		so(s, eq, "Andrew")
	})

	cv("GetInt", func() {
		i, err := o.GetInt("data", "year")
		so(err, isNil)
		so(i, eq, 2019)
	})

	cv("GetUint", func() {
		i, err := o.GetUint("data", "year")
		so(err, isNil)
		so(i, eq, 2019)
	})

	cv("GetInt64", func() {
		i, err := o.GetInt64("data", "year")
		so(err, isNil)
		so(i, eq, 2019)
	})

	cv("GetUint64", func() {
		i, err := o.GetUint64("data", "year")
		so(err, isNil)
		so(i, eq, 2019)
	})

	cv("GetInt32", func() {
		i, err := o.GetInt32("data", "negative")
		so(err, isNil)
		so(i, eq, -1234)
	})

	cv("Caseless.GetInt32", func() {
		i, err := o.Caseless().GetInt32("data", "negATive")
		so(err, isNil)
		so(i, eq, -1234)
	})

	cv("GetInt32 but not caseless", func() {
		_, err := o.GetInt32("data", "negATive")
		so(err, isErr)
	})

	cv("GetUint32", func() {
		i, err := o.GetUint32("data", "year")
		so(err, isNil)
		so(i, eq, 2019)
	})

	cv("GetFloat64", func() {
		f, err := o.GetFloat64("data", "YYYY.MM")
		so(err, isNil)
		so(f, eq, 2019.12)
	})

	cv("GetFloat32", func() {
		f, err := o.GetFloat32("data", "YYYY.MM")
		so(err, isNil)
		so(f, eq, 2019.12)
	})

	cv("GetNull", func() {
		err := o.GetNull("data", "message", -1)
		so(err, isNil)
	})

	cv("GetBool", func() {
		b, err := o.GetBool("data", "message", 2)
		so(err, isNil)
		so(b, isTrue)
	})

	cv("GetString from array of first one", func() {
		s, err := o.GetString("data", "message", 0)
		so(err, isNil)
		so(s, eq, "hello")
	})

	cv("GetString from array of last third one", func() {
		s, err := o.GetString("data", "message", -3)
		so(err, isNil)
		so(s, eq, "world")
	})

	cv("Len", func() {
		l := o.Len()
		so(l, eq, 1)

		v, _ := o.Get("data", "message")
		l = v.Len()
		so(l, eq, 4)

		v = o.MustGet("data", "message")
		l = v.Len()
		so(l, eq, 4)

		v, _ = o.Get("data", "author")
		l = v.Len()
		so(l, eq, 0)

		v = o.MustGet("data", "author")
		l = v.Len()
		so(l, eq, 0)
	})

	cv("GetObject", func() {
		v, err := o.GetObject("data")
		so(err, isNil)
		so(v.IsObject(), isTrue)
	})

	cv("GetObject in object", func() {
		v, err := o.Caseless().GetObject("Data")
		so(err, isNil)
		so(v.IsObject(), isTrue)

		v = o.Caseless().MustGet("Data")
		so(v.IsObject(), isTrue)
	})

	cv("nil V string", func() {
		v, _ := o.GetObject("not_exist")
		so(v.String(), eq, "")
	})

	cv("key: num_in_str", func() {
		v, err := o.Get("data", "num_in_str")
		so(err, isNil)
		so(v.Int(), eq, 2020)
		so(v.Float64(), eq, 2020.02)
		so(v.String(), eq, "2020.02")
	})

	cv("key: negative_in_str", func() {
		v, err := o.Get("data", "negative_in_str")
		so(err, isNil)
		so(v.Int(), eq, -12345)
		so(v.IsString(), isTrue)
		so(v.String(), eq, "-12345")
	})

	cv("key: invalid_num_in_str", func() {
		v, err := o.Get("data", "invalid_num_in_str")
		so(err, isNil)
		so(v.IsString(), isTrue)
		so(v.Float64(), eq, 0)
	})
}

func testMiscError(t *testing.T) {
	var err error
	raw := `{"array":[0,1,2,3],"string":"hello, world","number":1234.12345}`
	v, err := UnmarshalString(raw)
	so(err, isNil)

	// param error
	_, err = v.GetInt("array", true)
	so(err, isErr)
	_, err = v.GetString(true)
	so(err, isErr)

	// Caseless via non object or array
	child, err := v.Get("string")
	so(child, notNil)
	so(err, isNil)
	errV, err := child.Caseless().Get("NOT_EXIST")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)

	errV = child.Caseless().MustGet("NOT_EXIST")
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)

	// out of range
	errV, err = v.Get("array", 100)
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)

	errV = v.MustGet("array", 100)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)

	// not support
	errV, err = v.Get("string", "hello")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)

	// GetString
	s, err := v.GetString("number")
	so(err, isErr)
	so(s, eq, "")
	s, err = v.GetString("not exist")
	so(err, isErr)
	so(s, eq, "")

	// GetInt... and GetUint..
	i, err := v.GetInt("string")
	so(err, isErr)
	so(i, eq, 0)
	u, err := v.GetUint("string")
	so(err, isErr)
	so(u, eq, 0)
	i64, err := v.GetInt64("string")
	so(err, isErr)
	so(i64, eq, 0)
	u64, err := v.GetUint64("string")
	so(err, isErr)
	so(u64, eq, 0)
	i32, err := v.GetInt32("string")
	so(err, isErr)
	so(i32, eq, 0)
	u32, err := v.GetUint32("string")
	so(err, isErr)
	so(u32, eq, 0)
	f64, err := v.GetFloat64("string")
	so(err, isErr)
	so(f64, eq, 0.0)
	f32, err := v.GetFloat32("string")
	so(err, isErr)
	so(f32, eq, 0.0)

	// number not exist
	s, err = v.GetString("not exist")
	so(err, isErr)
	so(s, eq, "")
	i, err = v.GetInt("not exist")
	so(err, isErr)
	so(i, eq, 0)
	u, err = v.GetUint("not exist")
	so(err, isErr)
	so(u, eq, 0)
	i64, err = v.GetInt64("not exist")
	so(err, isErr)
	so(i64, eq, 0)
	u64, err = v.GetUint64("not exist")
	so(err, isErr)
	so(u64, eq, 0)
	i32, err = v.GetInt32("not exist")
	so(err, isErr)
	so(i32, isZero)
	u32, err = v.GetUint32("not exist")
	so(err, isErr)
	so(u32, isZero)
	f64, err = v.GetFloat64("not exist")
	so(err, isErr)
	so(f64, isZero)
	f32, err = v.GetFloat32("not exist")
	so(err, isErr)
	so(f32, isZero)

	// GetObject and GetArray
	errV, err = v.GetObject("string")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)
	errV, err = v.GetArray("string")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)
	errV, err = v.GetObject("not exist")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)
	errV, err = v.GetArray("not exist")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)

	// GetBool and GetNull
	bol, err := v.GetBool("string")
	so(err, isErr)
	so(bol, eq, false)
	err = v.GetNull("string")
	so(err, isErr)
	bol, err = v.GetBool("not exist")
	so(err, isErr)
	so(bol, eq, false)
	err = v.GetNull("not exist")
	so(err, isErr)

	// GetBytes
	byt, err := v.GetBytes("string")
	so(err, isErr)
	so(len(byt), eq, 0)
	byt, err = v.GetBytes("array")
	so(err, isErr)
	so(len(byt), eq, 0)
}

func testCaselessGet(t *testing.T) {
	raw := `{"data":{"STRING":"hello, world","INTEGER":12345,"TRUE":true,"FALSE":false,"NULL":null,"FLOAT":1234.5678,"OBJECT":{},"ARRAY":[]}}`

	v, err := UnmarshalString(raw)
	v.MustSetBytes([]byte{1, 2, 3, 4}).At("data", "BYTES")

	t.Log(v.MustMarshalString())

	so(err, isNil)
	so(v.IsObject(), isTrue)

	errV, err := v.Get("data", "object")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)
	_, err = v.Caseless().Get("data", "object")
	so(err, isNil)
	errV, err = v.Caseless().Get("data", "obj")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)

	errV, err = v.GetObject("data", "object")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)
	_, err = v.Caseless().GetObject("data", "object")
	so(err, isNil)

	errV, err = v.GetArray("data", "array")
	so(err, isErr)
	so(errV, notNil)
	so(errV.ValueType(), eq, NotExist)
	_, err = v.Caseless().GetArray("data", "array")
	so(err, isNil)

	_, err = v.GetBytes("data", "bytes")
	so(err, isErr)
	_, err = v.Caseless().GetBytes("data", "bytes")
	so(err, isNil)

	_, err = v.GetString("data", "string")
	so(err, isErr)
	_, err = v.Caseless().GetString("data", "string")
	so(err, isNil)

	_, err = v.GetInt("data", "integer")
	so(err, isErr)
	_, err = v.Caseless().GetInt("data", "integer")
	so(err, isNil)

	_, err = v.GetUint("data", "integer")
	so(err, isErr)
	_, err = v.Caseless().GetUint("data", "integer")
	so(err, isNil)

	_, err = v.GetInt64("data", "integer")
	so(err, isErr)
	_, err = v.Caseless().GetInt64("data", "integer")
	so(err, isNil)

	_, err = v.GetUint64("data", "integer")
	so(err, isErr)
	_, err = v.Caseless().GetUint64("data", "integer")
	so(err, isNil)

	_, err = v.GetInt32("data", "integer")
	so(err, isErr)
	_, err = v.Caseless().GetInt32("data", "integer")
	so(err, isNil)

	_, err = v.GetUint32("data", "integer")
	so(err, isErr)
	_, err = v.Caseless().GetUint32("data", "integer")
	so(err, isNil)

	_, err = v.GetFloat64("data", "float")
	so(err, isErr)
	_, err = v.Caseless().GetFloat64("data", "float")
	so(err, isNil)

	_, err = v.GetFloat32("data", "float")
	so(err, isErr)
	_, err = v.Caseless().GetFloat32("data", "float")
	so(err, isNil)

	_, err = v.GetBool("data", "true")
	so(err, isErr)
	_, err = v.Caseless().GetBool("data", "true")
	so(err, isNil)

	err = v.GetNull("data", "null")
	so(err, isErr)
	err = v.Caseless().GetNull("data", "null")
	so(err, isNil)

	err = v.Caseless().Delete("data", "array")
	so(err, isNil)
	sub, err := v.Caseless().Get("data", "array")
	so(err, isErr)
	so(sub, notNil)
	so(sub.ValueType(), eq, NotExist)
}

func testNotExistGet(t *testing.T) {
	cv("unmarshal a not exist V", func() {
		v := MustUnmarshalString("blahblah")
		so(v, notNil)
		so(v.ValueType(), eq, NotExist)

		sub, err := v.Get("string")
		so(sub, notNil)
		so(err, isErr)
		so(sub.ValueType(), eq, NotExist)
	})

	cv("not-exist-V.GetArray", func() {
		v, err := MustUnmarshalString("blahblah").GetArray("some_array", 1, 2, 3)
		so(v, notNil)
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)
	})

	cv("not-exist-V.GetObject", func() {
		v, err := MustUnmarshalString("blahblah").GetArray("some_object", 1, 2, 3)
		so(v, notNil)
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)
	})
}

func testGetNumFromString(t *testing.T) {
	cv("invalid number", func() {
		v := MustUnmarshalString(`{"num":"abcde","bool":true}`)

		i, err := v.GetInt("num")
		so(err, isErr)
		so(errors.Is(err, ErrParseNumberFromString), isTrue)
		so(i, isZero)

		f, err := v.GetFloat64("num")
		so(err, isErr)
		so(errors.Is(err, ErrParseNumberFromString), isTrue)
		so(f, isZero)

		f, err = v.GetFloat64("bool")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(f, isZero)
	})

	cv("int", func() {
		v, err := UnmarshalString(`{"num":"-123.25","negative":-9223372036854775808}`)
		so(err, isNil)

		i, err := v.GetInt("num")
		so(err, isErr)
		so(i, eq, -123)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)

		f, err := v.GetFloat64("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(f, eq, -123.25)

		i, err = v.GetInt(`negative`)
		so(err, isNil)
		so(i, eq, -9223372036854775808)
	})

	cv("uint", func() {
		v := MustUnmarshalString(`{"num":"123.25"}`)

		i, err := v.GetInt("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(i, eq, 123)

		f, err := v.GetFloat64("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(f, eq, 123.25)
	})

	cv("int32", func() {
		v := MustUnmarshalString(`{"num":"-123.25"}`)

		i, err := v.GetInt32("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(i, eq, -123)

		f, err := v.GetFloat32("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(f, eq, -123.25)
	})

	cv("uint32", func() {
		v := MustUnmarshalString(`{"num":"123.25"}`)

		i, err := v.GetInt("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(i, eq, 123)

		f, err := v.GetFloat32("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(f, eq, 123.25)
	})

	cv("int64", func() {
		v := MustUnmarshalString(`{"num":"-123.25"}`)

		i, err := v.GetInt64("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(i, eq, -123)

		f, err := v.GetFloat64("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(f, eq, -123.25)
	})

	cv("uint64", func() {
		v := MustUnmarshalString(`{"num":"123.25"}`)

		i, err := v.GetInt64("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(i, eq, 123)

		f, err := v.GetFloat64("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(f, eq, 123.25)
	})

	cv("bool", func() {
		v := MustUnmarshalString(`{"str":"true"}`)
		b, err := v.GetBool("str")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(b, isTrue)

		v = MustUnmarshalString(`{"str":""}`)
		b, err = v.GetBool("str")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(b, isFalse)

		v = MustUnmarshalString(`{"num":"0"}`)
		b, err = v.GetBool("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(b, isFalse)

		v = MustUnmarshalString(`{"num":1}`)
		b, err = v.GetBool("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(b, isTrue)

		v = MustUnmarshalString(`{"num":0}`)
		b, err = v.GetBool("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(b, isFalse)

		v = MustUnmarshalString(`{"num":null}`)
		b, err = v.GetBool("num")
		so(err, isErr)
		so(errors.Is(err, ErrTypeNotMatch), isTrue)
		so(b, isFalse)
	})
}
