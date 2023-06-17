package jsonvalue

import (
	"testing"
)

func testNewXxx(t *testing.T) {
	cv("NewString", func() { testNewString(t) })
	cv("NewBool", func() { testNewBool(t) })
	cv("NewNull", func() { testNewNull(t) })
	cv("NewIntXxx/UintXxx", func() { testNewInteger(t) })
	cv("NewFloat64/32", func() { testNewFloat(t) })
	cv("empty object/array", func() { testEmptyObjectArray(t) })
	cv("misc value", func() { testMiscValue(t) })
	cv("MustMarshal error", func() { testMustMarshalError(t) })
	cv("ValueError error", func() { testValueError(t) })
}

func testNewString(t *testing.T) {
	s := "你好，世界"
	v := NewString(s)
	so(v.String(), eq, s)
	so(v.ValueType(), eq, String)
}

func testNewBool(t *testing.T) {
	v := NewBool(true)
	so(v.Bool(), isTrue)
	so(v.IsBoolean(), isTrue)
	so(v.ValueType(), eq, Boolean)

	v = NewBool(false)
	so(v.Bool(), isFalse)
	so(v.IsBoolean(), isTrue)
	so(v.ValueType(), eq, Boolean)
}

func testNewNull(t *testing.T) {
	v := NewNull()
	so(v.IsNull(), isTrue)
	so(v.ValueType(), eq, Null)
}

func testNewInteger(t *testing.T) {
	i := int64(-1234567)

	v := NewInt(int(i))
	so(v.Int(), eq, int(i))
	so(v.ValueType(), eq, Number)

	v = NewUint(uint(i))
	so(v.Uint(), eq, uint(i))
	so(v.ValueType(), eq, Number)

	v = NewInt32(int32(i))
	so(v.Int32(), eq, int32(i))
	so(v.ValueType(), eq, Number)

	v = NewUint32(uint32(i))
	so(v.Uint32(), eq, uint32(i))
	so(v.ValueType(), eq, Number)

	v = NewInt64(int64(i))
	so(v.Int64(), eq, int64(i))
	so(v.ValueType(), eq, Number)

	v = NewUint64(uint64(i))
	so(v.Uint64(), eq, uint64(i))
	so(v.ValueType(), eq, Number)
}

func testNewFloat(t *testing.T) {
	s := "3.1415926535"
	f := 3.1415926535

	v := NewFloat64f(f, 'f', 10)
	so(v.String(), eq, s)
	so(v.ValueType(), eq, Number)

	v = NewFloat64f(f, '?', 11)
	so(v.String(), eq, s)
	so(v.ValueType(), eq, Number)

	v = NewFloat64f(f, 'g', 3)
	so(v.String(), eq, "3.14")
	so(v.ValueType(), eq, Number)

	s = "3.1415927"
	v = NewFloat32(float32(f))
	so(v.String(), eq, s)
	so(v.ValueType(), eq, Number)

	v = NewFloat32f(float32(f), 'f', 5)
	so(v.String(), eq, "3.14159")
	so(v.ValueType(), eq, Number)

	v = NewFloat32f(float32(f), 'e', 5)
	so(v.String(), eq, "3.14159e+00")
	so(v.ValueType(), eq, Number)

	v = NewFloat32f(float32(f), '?', 5)
	so(v.String(), eq, "3.1416")
	so(v.ValueType(), eq, Number)
}

func testEmptyObjectArray(t *testing.T) {
	v := NewObject()
	b, _ := v.Marshal()
	so(string(b), eq, "{}")
	so(v.ValueType(), eq, Object)

	v = NewArray()
	b = v.MustMarshal()
	so(string(b), eq, "[]")
	so(v.ValueType(), eq, Array)
}

func testMiscValue(t *testing.T) {
	cv("parse array", func() {
		raw := "\r\n[1, 2, 3 ]\t\b"
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsArray(), isTrue)
	})

	cv("parse object", func() {
		raw := `{ }`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)
	})

	cv("parse string", func() {
		raw := ` "hello, world"  `
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsString(), isTrue)
		so(v.Int(), isZero)
		so(v.Uint(), isZero)
		so(v.Int64(), isZero)
		so(v.Uint64(), isZero)
		so(v.Int32(), isZero)
		so(v.Uint32(), isZero)
		so(v.Float64(), isZero)
		so(v.Float32(), isZero)
		so(v.IsFloat(), isFalse)
		so(v.IsInteger(), isFalse)
	})

	cv("parse string with special character", func() {
		raw := `"\"\\\/\f\t\r\n\b\u0030\uD87E\uDC04"`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsString(), isTrue)
	})

	cv("parse null", func() {
		raw := `null`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsNull(), isTrue)
	})

	cv("parse bool (true)", func() {
		raw := `true`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsBoolean(), isTrue)
		so(v.Bool(), isTrue)
	})

	cv("parse bool (false)", func() {
		raw := `false`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsBoolean(), isTrue)
		so(v.Bool(), isFalse)
	})

	cv("parse negative float", func() {
		raw := `-12345.12345`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsNumber(), isTrue)
		so(v.IsFloat(), isTrue)
		so(v.IsNegative(), isTrue)
		so(v.GreaterThanInt64Max(), isFalse)
	})

	cv("parse exponential form float", func() {
		raw := `-1.25e3`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsFloat(), isTrue)
		so(v.IsNegative(), isTrue)
		so(v.Float32(), eq, -1250)

		raw = `1.25E-1`
		v, err = UnmarshalString(raw)
		so(err, isNil)
		so(v.IsFloat(), isTrue)
		so(v.IsNegative(), isFalse)
		so(v.Float32(), eq, 0.125)
	})

	cv("parse negative integer", func() {
		raw := `-12345`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsNumber(), isTrue)
		so(v.IsFloat(), isFalse)
		so(v.IsNegative(), isTrue)
		so(v.GreaterThanInt64Max(), isFalse)
	})

	cv("parse positive integer", func() {
		raw := `12345`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsNumber(), isTrue)
		so(v.IsFloat(), isFalse)
		so(v.IsPositive(), isTrue)
		so(v.GreaterThanInt64Max(), isFalse)
	})

	cv("parse big uint64", func() {
		raw := `18446744073709551615` // 0xFFFFFFFFFFFFFFFF
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsNumber(), isTrue)
		so(v.IsFloat(), isFalse)
		so(v.IsPositive(), isTrue)
		so(v.GreaterThanInt64Max(), isTrue)
	})

	cv("parse object in array", func() {
		raw := `[{}]`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsArray(), isTrue)

		c, err := v.Get(0)
		so(err, isNil)
		so(c.IsObject(), isTrue)
	})

	cv("parse array in object", func() {
		raw := `{"array":[]}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("array")
		so(err, isNil)
		so(c.IsArray(), isTrue)
	})

	cv("parse float in object", func() {
		raw := `{"float":123.4567}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("float")
		so(err, isNil)
		so(c.IsFloat(), isTrue)
	})

	cv("parse integer in object", func() {
		raw := `{"int":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("int")
		so(err, isNil)
		so(c.IsInteger(), isTrue)
	})

	cv("parse int in object", func() {
		raw := `{"int":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("int")
		so(err, isNil)
		so(c.Int(), eq, 123)
	})

	cv("parse uint in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("uint")
		so(err, isNil)
		so(c.Uint(), eq, 123)
	})

	cv("parse int64 in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("uint")
		so(err, isNil)
		so(c.Int64(), eq, 123)
	})

	cv("parse uint64 in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("uint")
		so(err, isNil)
		so(c.Uint64(), eq, 123)
	})

	cv("parse int32 in object", func() {
		raw := `{"int":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("int")
		so(err, isNil)
		so(c.Int32(), eq, 123)
	})

	cv("parse uint32 in object", func() {
		raw := `{"uint":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("uint")
		so(err, isNil)
		so(c.Uint32(), eq, 123)
	})

	cv("parse float32 in object", func() {
		raw := `{"float":123.456}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("float")
		so(err, isNil)
		so(c.Float32(), eq, 123.456)
	})

	cv("parse float64 in object", func() {
		raw := `{"float":123.456}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("float")
		so(err, isNil)
		so(c.Float64(), eq, 123.456)
	})

	cv("parse negative in object", func() {
		raw := `{"negative":-123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("negative")
		so(err, isNil)
		so(c.IsNegative(), isTrue)
	})

	cv("parse positive in object", func() {
		raw := `{"positive":123}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("positive")
		so(err, isNil)
		so(c.IsPositive(), isTrue)
	})

	cv("parse greater than int64 in object", func() {
		raw := `{"int":9223372036854775808}`
		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		c, err := v.Get("int")
		so(err, isNil)
		so(c.GreaterThanInt64Max(), isTrue)
	})

	cv("create an initialized JSON object", func() {
		v := NewObject(map[string]any{
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

		so(v.GetNull("null"), isNil)

		str, err := v.GetString("string")
		so(err, isNil)
		so(str, eq, "string")

		bl, err := v.GetBool("true")
		so(err, isNil)
		so(bl, isTrue)

		checkInt := func(k string, target int) {
			i, err := v.GetInt(k)
			so(err, isNil)
			so(i, eq, target)
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
		so(err, isNil)
		so(f32, eq, float32(-1.1))

		f64, err := v.GetFloat64("float64")
		so(err, isNil)
		so(f64, eq, float64(-2.2))
	})
}

func testMustMarshalError(t *testing.T) {
	shouldPanic(func() {
		v := &V{}
		v.MustMarshal()
	})

	shouldPanic(func() {
		v := &V{}
		v.MustMarshalString()
	})
}

func testValueError(t *testing.T) {
	var err error
	var raw string
	var v *V

	cv("invalid jsonvalue.V", func() {
		v = &V{}
		so(v.String(), eq, "")
	})

	cv("number in string", func() {
		v, err = UnmarshalString(`"12.ABCD"`)
		so(err, isNil)
		so(v.Int(), isZero)

		v, err = UnmarshalString(`"-12ABCD"`)
		so(err, isNil)
		so(v.Int(), isZero)

		v, err = UnmarshalString(`{}`)
		so(err, isNil)
		so(v.Int(), isZero)

		v, err = UnmarshalString(`"1234"`)
		so(err, isNil)
		so(v.Int(), eq, 1234)
		so(v.IsNumber(), isFalse)
		so(v.IsFloat(), isFalse)
		so(v.IsInteger(), isFalse)
		so(v.IsNegative(), isFalse)
		so(v.IsPositive(), isFalse)
		so(v.GreaterThanInt64Max(), isFalse)
	})

	cv("nil string input", func() {
		_, err = UnmarshalString("")
		so(err, isErr)
	})

	cv("nil bytes input", func() {
		_, err = Unmarshal(nil)
		so(err, isErr)
	})

	cv("illegal char", func() {
		_, err = UnmarshalString(`\\`)
		so(err, isErr)
	})

	cv("no start chars", func() {
		_, err = UnmarshalString(`      `)
		so(err, isErr)
	})

	cv("illegal float number", func() {
		_, err = UnmarshalString(`1.a`)
		so(err, isErr)

		_, err = UnmarshalString(`1.`)
		so(err, isErr)

		_, err = UnmarshalString(`.1`)
		so(err, isErr)
	})

	cv("illegal negative interger", func() {
		_, err = UnmarshalString(`-1a`)
		so(err, isErr)
	})

	cv("illegal positive interger", func() {
		_, err = UnmarshalString(`11a`)
		so(err, isErr)
	})

	cv("illegal true", func() {
		_, err = UnmarshalString(`trUE`)
		so(err, isErr)
	})

	cv("illegal false", func() {
		_, err = UnmarshalString(`fAlse`)
		so(err, isErr)
	})

	cv("illegal null", func() {
		_, err = UnmarshalString(`nUll`)
		so(err, isErr)
	})

	cv("illegal string", func() {
		_, err = UnmarshalString(`"too many quote""`)
		so(err, isErr)
	})

	cv("too short string", func() {
		_, err = UnmarshalString(`"`)
		so(err, isErr)
	})

	cv("illegal escaping", func() {
		_, err = UnmarshalString(`"\"`)
		so(err, isErr)
	})

	cv("illegal bool in object", func() {
		_, err = UnmarshalString(`{"bool":tRue}`)
		so(err, isErr)
	})

	cv("illegal bool in array", func() {
		_, err = UnmarshalString(`[tRue]`)
		so(err, isErr)
	})

	cv("illegal array", func() {
		_, err = UnmarshalString(`["incompleteString]`)
		so(err, isErr)
	})

	cv("illegal array without ]", func() {
		_, err = UnmarshalString(`[   `)
		so(err, isErr)
	})

	cv("illegal object in array without }", func() {
		_, err = UnmarshalString(`[{   ]`)
		so(err, isErr)
	})

	cv("illegal array in array without ]", func() {
		_, err = UnmarshalString(`[[  224 ]`)
		so(err, isErr)
	})

	cv("another illegal array in array without ]", func() {
		_, err = UnmarshalString(`[ [  `)
		so(err, isErr)
	})

	cv("illegal number in array without ]", func() {
		_, err = UnmarshalString(`[224.. ]`)
		so(err, isErr)
	})

	cv("illegal number in array without decimal part", func() {
		_, err = UnmarshalString(`[224.]`)
		so(err, isErr)
	})

	cv("another illegal number in array without ]", func() {
		_, err = UnmarshalString(`[-18446744073709551615 ]`)
		so(err, isErr)
	})

	cv("illegal false in array ", func() {
		_, err = UnmarshalString(`[fASLE ]`)
		so(err, isErr)
	})

	cv("illegal true in array ", func() {
		_, err = UnmarshalString(`[tRUE ]`)
		so(err, isErr)
	})

	cv("illegal null in array ", func() {
		_, err = UnmarshalString(`[nULL ]`)
		so(err, isErr)
	})

	cv("illegal character in array ", func() {
		_, err = UnmarshalString(`[W]`)
		so(err, isErr)
	})

	cv("marshaling uninitialized value", func() {
		v = &V{}
		_, err = v.MarshalString()
		so(err, isErr)

		_, err = v.Marshal()
		so(err, isErr)
	})

	cv("marshaling with option", func() {
		v = NewObject()
		v.MustSetNull().At("null")
		raw, _ = v.MarshalString()
		so(raw, eq, `{"null":null}`)

		raw, _ = v.MarshalString(Opt{OmitNull: true})
		so(raw, eq, `{}`)

		raw, _ = v.MarshalString(OptOmitNull(true))
		so(raw, eq, `{}`)

		rawB, _ := v.Marshal(Opt{OmitNull: true})
		so(string(rawB), eq, `{}`)

		rawB, _ = v.Marshal(OptOmitNull(true))
		so(string(rawB), eq, `{}`)
	})

	cv("illegal kvs in object", func() {
		_, err = UnmarshalString(`{true}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{:"value"}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{"key" true}`) // missing colon
		so(err, isErr)

		_, err = UnmarshalString(`{"key":}`) // missing value
		so(err, isErr)

		_, err = UnmarshalString(`{"key":"value"   `) // missing }
		so(err, isErr)

		_, err = UnmarshalString(`{"key":"value"`) // missing }
		so(err, isErr)

		_, err = UnmarshalString(`{"key":,}`) // missing value
		so(err, isErr)

		_, err = UnmarshalString(`{"key"::"value"}`) // duplicate colon
		so(err, isErr)

		_, err = UnmarshalString(`{{}}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{"object":{ILLEGAL}}`) // invalid object in object
		so(err, isErr)

		_, err = UnmarshalString(`{"array":[ILLEGAL]}`) // invalid array in object
		so(err, isErr)

		_, err = UnmarshalString(`{[]}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{[}`) // missing ]
		so(err, isErr)

		_, err = UnmarshalString(`{12345}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{"big_int":-18446744073709551615}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{"key" "value"}`) // missing colon
		so(err, isErr)

		_, err = UnmarshalString(`{"key":"\"}`) // invalid string in object
		so(err, isErr)

		_, err = UnmarshalString(`{"key\u":"value"}`) // invalid key in object
		so(err, isErr)

		_, err = UnmarshalString(`{false}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{"false":fAlse}`) // illegal value
		so(err, isErr)

		_, err = UnmarshalString(`{null}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{"null":nUll}`) // illegal value
		so(err, isErr)

		_, err = UnmarshalString(`{true}`) // missing key
		so(err, isErr)

		_, err = UnmarshalString(`{"true":tRue}`) // illegal value
		so(err, isErr)
	})
}
