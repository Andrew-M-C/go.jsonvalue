package jsonvalue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewXxx(t *testing.T) {
	test(t, "NewString", testNewString)
	test(t, "NewBool", testNewBool)
	test(t, "NewNull", testNewNull)
	test(t, "NewIntXxx/UintXxx", testNewInteger)
	test(t, "NewFloat64/32", testNewFloat)
	test(t, "empty object/array", testEmptyObjectArray)
	test(t, "misc value", testMiscValue)
	test(t, "MustMarshal error", testMustMarshalError)
	test(t, "ValueError error", testValueError)
}

func testNewString(t *testing.T) {
	s := "你好，世界"
	v := NewString(s)
	So(v.String(), ShouldEqual, s)
	So(v.ValueType(), ShouldEqual, String)
}

func testNewBool(t *testing.T) {
	v := NewBool(true)
	So(v.Bool(), ShouldBeTrue)
	So(v.IsBoolean(), ShouldBeTrue)
	So(v.ValueType(), ShouldEqual, Boolean)

	v = NewBool(false)
	So(v.Bool(), ShouldBeFalse)
	So(v.IsBoolean(), ShouldBeTrue)
	So(v.ValueType(), ShouldEqual, Boolean)
}

func testNewNull(t *testing.T) {
	v := NewNull()
	So(v.IsNull(), ShouldBeTrue)
	So(v.ValueType(), ShouldEqual, Null)
}

func testNewInteger(t *testing.T) {
	i := int64(-1234567)

	v := NewInt(int(i))
	So(v.Int(), ShouldEqual, int(i))
	So(v.ValueType(), ShouldEqual, Number)

	v = NewUint(uint(i))
	So(v.Uint(), ShouldEqual, uint(i))
	So(v.ValueType(), ShouldEqual, Number)

	v = NewInt32(int32(i))
	So(v.Int32(), ShouldEqual, int32(i))
	So(v.ValueType(), ShouldEqual, Number)

	v = NewUint32(uint32(i))
	So(v.Uint32(), ShouldEqual, uint32(i))
	So(v.ValueType(), ShouldEqual, Number)

	v = NewInt64(int64(i))
	So(v.Int64(), ShouldEqual, int64(i))
	So(v.ValueType(), ShouldEqual, Number)

	v = NewUint64(uint64(i))
	So(v.Uint64(), ShouldEqual, uint64(i))
	So(v.ValueType(), ShouldEqual, Number)
}

func testNewFloat(t *testing.T) {
	s := "3.1415926535"
	f := 3.1415926535

	v := NewFloat64(f, 10)
	So(v.String(), ShouldEqual, s)
	So(v.ValueType(), ShouldEqual, Number)

	v = NewFloat64(f, 2)
	So(v.String(), ShouldEqual, "3.14")
	So(v.ValueType(), ShouldEqual, Number)

	s = "3.1415927"
	v = NewFloat32(float32(f), -1)
	So(v.String(), ShouldEqual, s)
	So(v.ValueType(), ShouldEqual, Number)

	v = NewFloat32(float32(f), 5)
	So(v.String(), ShouldEqual, "3.14159")
	So(v.ValueType(), ShouldEqual, Number)
}

func testEmptyObjectArray(t *testing.T) {
	v := NewObject()
	b, _ := v.Marshal()
	So(string(b), ShouldEqual, "{}")
	So(v.ValueType(), ShouldEqual, Object)

	v = NewArray()
	b = v.MustMarshal()
	So(string(b), ShouldEqual, "[]")
	So(v.ValueType(), ShouldEqual, Array)
}

func testMiscValue(t *testing.T) {
	Convey("parse array", func() {
		raw := "\r\n[1, 2, 3 ]\t\b"
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsArray(), ShouldBeTrue)
	})

	Convey("parse object", func() {
		raw := `{ }`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)
	})

	Convey("parse string", func() {
		raw := ` "hello, world"  `
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsString(), ShouldBeTrue)
		So(v.Int(), ShouldBeZeroValue)
		So(v.Uint(), ShouldBeZeroValue)
		So(v.Int64(), ShouldBeZeroValue)
		So(v.Uint64(), ShouldBeZeroValue)
		So(v.Int32(), ShouldBeZeroValue)
		So(v.Uint32(), ShouldBeZeroValue)
		So(v.Float64(), ShouldBeZeroValue)
		So(v.Float32(), ShouldBeZeroValue)
		So(v.IsFloat(), ShouldBeFalse)
		So(v.IsInteger(), ShouldBeFalse)
	})

	Convey("parse string with special character", func() {
		raw := `"\"\\\/\f\t\r\n\b\u0030\uD87E\uDC04"`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsString(), ShouldBeTrue)
	})

	Convey("parse null", func() {
		raw := `null`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsNull(), ShouldBeTrue)
	})

	Convey("parse bool (true)", func() {
		raw := `true`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsBoolean(), ShouldBeTrue)
		So(v.Bool(), ShouldBeTrue)
	})

	Convey("parse bool (false)", func() {
		raw := `false`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsBoolean(), ShouldBeTrue)
		So(v.Bool(), ShouldBeFalse)
	})

	Convey("parse negative float", func() {
		raw := `-12345.12345`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.IsFloat(), ShouldBeTrue)
		So(v.IsNegative(), ShouldBeTrue)
		So(v.GreaterThanInt64Max(), ShouldBeFalse)
	})

	Convey("parse exponential form float", func() {
		raw := `-1.25e3`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsFloat(), ShouldBeTrue)
		So(v.IsNegative(), ShouldBeTrue)
		So(v.Float32(), ShouldEqual, -1250)

		raw = `1.25E-1`
		v, err = UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsFloat(), ShouldBeTrue)
		So(v.IsNegative(), ShouldBeFalse)
		So(v.Float32(), ShouldEqual, 0.125)
	})

	Convey("parse negative integer", func() {
		raw := `-12345`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.IsFloat(), ShouldBeFalse)
		So(v.IsNegative(), ShouldBeTrue)
		So(v.GreaterThanInt64Max(), ShouldBeFalse)
	})

	Convey("parse positive integer", func() {
		raw := `12345`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.IsFloat(), ShouldBeFalse)
		So(v.IsPositive(), ShouldBeTrue)
		So(v.GreaterThanInt64Max(), ShouldBeFalse)
	})

	Convey("parse big uint64", func() {
		raw := `18446744073709551615` // 0xFFFFFFFFFFFFFFFF
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsNumber(), ShouldBeTrue)
		So(v.IsFloat(), ShouldBeFalse)
		So(v.IsPositive(), ShouldBeTrue)
		So(v.GreaterThanInt64Max(), ShouldBeTrue)
	})

	Convey("parse object in array", func() {
		raw := `[{}]`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsArray(), ShouldBeTrue)

		c, err := v.Get(0)
		So(err, ShouldBeNil)
		So(c.IsObject(), ShouldBeTrue)
	})

	Convey("parse array in object", func() {
		raw := `{"array":[]}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("array")
		So(err, ShouldBeNil)
		So(c.IsArray(), ShouldBeTrue)
	})

	Convey("parse float in object", func() {
		raw := `{"float":123.4567}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("float")
		So(err, ShouldBeNil)
		So(c.IsFloat(), ShouldBeTrue)
	})

	Convey("parse integer in object", func() {
		raw := `{"int":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("int")
		So(err, ShouldBeNil)
		So(c.IsInteger(), ShouldBeTrue)
	})

	Convey("parse int in object", func() {
		raw := `{"int":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("int")
		So(err, ShouldBeNil)
		So(c.Int(), ShouldEqual, 123)
	})

	Convey("parse uint in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("uint")
		So(err, ShouldBeNil)
		So(c.Uint(), ShouldEqual, 123)
	})

	Convey("parse int64 in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("uint")
		So(err, ShouldBeNil)
		So(c.Int64(), ShouldEqual, 123)
	})

	Convey("parse uint64 in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("uint")
		So(err, ShouldBeNil)
		So(c.Uint64(), ShouldEqual, 123)
	})

	Convey("parse int32 in object", func() {
		raw := `{"int":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("int")
		So(err, ShouldBeNil)
		So(c.Int32(), ShouldEqual, 123)
	})

	Convey("parse uint32 in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("uint")
		So(err, ShouldBeNil)
		So(c.Uint32(), ShouldEqual, 123)
	})

	Convey("parse float32 in object", func() {
		raw := `{"float":123.456}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("float")
		So(err, ShouldBeNil)
		So(c.Float32(), ShouldEqual, 123.456)
	})

	Convey("parse float64 in object", func() {
		raw := `{"float":123.456}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("float")
		So(err, ShouldBeNil)
		So(c.Float64(), ShouldEqual, 123.456)
	})

	Convey("parse negative in object", func() {
		raw := `{"negative":-123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("negative")
		So(err, ShouldBeNil)
		So(c.IsNegative(), ShouldBeTrue)
	})

	Convey("parse positive in object", func() {
		raw := `{"positive":123}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("positive")
		So(err, ShouldBeNil)
		So(c.IsPositive(), ShouldBeTrue)
	})

	Convey("parse greater than int64 in object", func() {
		raw := `{"int":9223372036854775808}`
		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		c, err := v.Get("int")
		So(err, ShouldBeNil)
		So(c.GreaterThanInt64Max(), ShouldBeTrue)
	})

	Convey("create an initialized JSON object", func() {
		v := NewObject(map[string]interface{}{
			"null":    nil,
			"string":  "string",
			"true":    true,
			"int":     int(-1),
			"uint":    uint(1),
			"int8":    int8(-2),
			"uint8":   uint8(2),
			"int16":   int16(-3),
			"uint16":  uint16(3),
			"int32":   int32(-4),
			"uint32":  uint32(4),
			"int64":   int64(-5),
			"uint64":  uint64(5),
			"float32": float32(-1.1),
			"float64": float64(-2.2),
		})

		So(v.GetNull("null"), ShouldBeNil)

		str, err := v.GetString("string")
		So(err, ShouldBeNil)
		So(str, ShouldEqual, "string")

		bl, err := v.GetBool("true")
		So(err, ShouldBeNil)
		So(bl, ShouldBeTrue)

		checkInt := func(k string, target int) {
			i, err := v.GetInt(k)
			So(err, ShouldBeNil)
			So(i, ShouldEqual, target)
		}
		checkInt("int", -1)
		checkInt("uint", 1)
		checkInt("int8", -2)
		checkInt("uint8", 2)
		checkInt("int16", -3)
		checkInt("uint16", 3)
		checkInt("int32", -4)
		checkInt("uint32", 4)
		checkInt("int64", -5)
		checkInt("uint64", 5)

		f32, err := v.GetFloat32("float32")
		So(err, ShouldBeNil)
		So(f32, ShouldEqual, float32(-1.1))

		f64, err := v.GetFloat64("float64")
		So(err, ShouldBeNil)
		So(f64, ShouldEqual, float64(-2.2))
	})
}

func testMustMarshalError(t *testing.T) {
	ShouldPanic(func() {
		v := &V{}
		v.MustMarshal()
	})

	ShouldPanic(func() {
		v := &V{}
		v.MustMarshalString()
	})
}

func testValueError(t *testing.T) {
	var err error
	var raw string
	var v *V

	Convey("invalid jsonvalue.V", func() {
		v = &V{}
		So(v.String(), ShouldEqual, "")
	})

	Convey("number in string", func() {
		v, err = UnmarshalString(`"12.ABCD"`)
		So(err, ShouldBeNil)
		So(v.Int(), ShouldBeZeroValue)

		v, err = UnmarshalString(`"-12ABCD"`)
		So(err, ShouldBeNil)
		So(v.Int(), ShouldBeZeroValue)

		v, err = UnmarshalString(`{}`)
		So(err, ShouldBeNil)
		So(v.Int(), ShouldBeZeroValue)

		v, err = UnmarshalString(`"1234"`)
		So(err, ShouldBeNil)
		So(v.Int(), ShouldEqual, 1234)
		So(v.IsNumber(), ShouldBeFalse)
		So(v.IsFloat(), ShouldBeFalse)
		So(v.IsInteger(), ShouldBeFalse)
		So(v.IsNegative(), ShouldBeFalse)
		So(v.IsPositive(), ShouldBeFalse)
		So(v.GreaterThanInt64Max(), ShouldBeFalse)
	})

	Convey("nil string input", func() {
		_, err = UnmarshalString("")
		So(err, ShouldBeError)
	})

	Convey("nil bytes input", func() {
		_, err = Unmarshal(nil)
		So(err, ShouldBeError)
	})

	Convey("illegal char", func() {
		_, err = UnmarshalString(`\\`)
		So(err, ShouldBeError)
	})

	Convey("no start chars", func() {
		_, err = UnmarshalString(`      `)
		So(err, ShouldBeError)
	})

	Convey("illegal float number", func() {
		_, err = UnmarshalString(`1.a`)
		So(err, ShouldBeError)

		_, err = UnmarshalString(`1.`)
		So(err, ShouldBeError)

		_, err = UnmarshalString(`.1`)
		So(err, ShouldBeError)
	})

	Convey("illegal negative interger", func() {
		_, err = UnmarshalString(`-1a`)
		So(err, ShouldBeError)
	})

	Convey("illegal positive interger", func() {
		_, err = UnmarshalString(`11a`)
		So(err, ShouldBeError)
	})

	Convey("illegal true", func() {
		_, err = UnmarshalString(`trUE`)
		So(err, ShouldBeError)
	})

	Convey("illegal false", func() {
		_, err = UnmarshalString(`fAlse`)
		So(err, ShouldBeError)
	})

	Convey("illegal null", func() {
		_, err = UnmarshalString(`nUll`)
		So(err, ShouldBeError)
	})

	Convey("illegal string", func() {
		_, err = UnmarshalString(`"too many quote""`)
		So(err, ShouldBeError)
	})

	Convey("too short string", func() {
		_, err = UnmarshalString(`"`)
		So(err, ShouldBeError)
	})

	Convey("illegal escaping", func() {
		_, err = UnmarshalString(`"\"`)
		So(err, ShouldBeError)
	})

	Convey("illegal bool in object", func() {
		_, err = UnmarshalString(`{"bool":tRue}`)
		So(err, ShouldBeError)
	})

	Convey("illegal bool in array", func() {
		_, err = UnmarshalString(`[tRue]`)
		So(err, ShouldBeError)
	})

	Convey("illegal array", func() {
		_, err = UnmarshalString(`["incompleteString]`)
		So(err, ShouldBeError)
	})

	Convey("illegal array without ]", func() {
		_, err = UnmarshalString(`[   `)
		So(err, ShouldBeError)
	})

	Convey("illegal object in array without }", func() {
		_, err = UnmarshalString(`[{   ]`)
		So(err, ShouldBeError)
	})

	Convey("illegal array in array without ]", func() {
		_, err = UnmarshalString(`[[  224 ]`)
		So(err, ShouldBeError)
	})

	Convey("another illegal array in array without ]", func() {
		_, err = UnmarshalString(`[ [  `)
		So(err, ShouldBeError)
	})

	Convey("illegal number in array without ]", func() {
		_, err = UnmarshalString(`[224.. ]`)
		So(err, ShouldBeError)
	})

	Convey("illegal number in array without decimal part", func() {
		_, err = UnmarshalString(`[224.]`)
		So(err, ShouldBeError)
	})

	Convey("another illegal number in array without ]", func() {
		_, err = UnmarshalString(`[-18446744073709551615 ]`)
		So(err, ShouldBeError)
	})

	Convey("illegal false in array ", func() {
		_, err = UnmarshalString(`[fASLE ]`)
		So(err, ShouldBeError)
	})

	Convey("illegal true in array ", func() {
		_, err = UnmarshalString(`[tRUE ]`)
		So(err, ShouldBeError)
	})

	Convey("illegal null in array ", func() {
		_, err = UnmarshalString(`[nULL ]`)
		So(err, ShouldBeError)
	})

	Convey("illegal character in array ", func() {
		_, err = UnmarshalString(`[W]`)
		So(err, ShouldBeError)
	})

	Convey("marshaling uninitialized value", func() {
		v = &V{}
		_, err = v.MarshalString()
		So(err, ShouldBeError)

		_, err = v.Marshal()
		So(err, ShouldBeError)
	})

	Convey("marshaling with option", func() {
		v = NewObject()
		v.SetNull().At("null")
		raw, _ = v.MarshalString()
		So(raw, ShouldEqual, `{"null":null}`)

		raw, _ = v.MarshalString(Opt{OmitNull: true})
		So(raw, ShouldEqual, `{}`)

		raw, _ = v.MarshalString(OptOmitNull(true))
		So(raw, ShouldEqual, `{}`)

		rawB, _ := v.Marshal(Opt{OmitNull: true})
		So(string(rawB), ShouldEqual, `{}`)

		rawB, _ = v.Marshal(OptOmitNull(true))
		So(string(rawB), ShouldEqual, `{}`)
	})

	Convey("illegal kvs in object", func() {
		_, err = UnmarshalString(`{true}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{:"value"}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key" true}`) // missing colon
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key":}`) // missing value
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key":"value"   `) // missing }
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key":"value"`) // missing }
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key":,}`) // missing value
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key"::"value"}`) // duplicate colon
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{{}}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"object":{ILLEGAL}}`) // invalid object in object
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"array":[ILLEGAL]}`) // invalid array in object
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{[]}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{[}`) // missing ]
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{12345}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"big_int":-18446744073709551615}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key" "value"}`) // missing colon
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key":"\"}`) // invalid string in object
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"key\u":"value"}`) // invalid key in object
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{false}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"false":fAlse}`) // illegal value
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{null}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"null":nUll}`) // illegal value
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{true}`) // missing key
		So(err, ShouldBeError)

		_, err = UnmarshalString(`{"true":tRue}`) // illegal value
		So(err, ShouldBeError)
	})
}
