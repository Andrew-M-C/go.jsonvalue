package jsonvalue

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So

	eq = convey.ShouldEqual
	ne = convey.ShouldNotEqual

	isNil   = convey.ShouldBeNil
	notNil  = convey.ShouldNotBeNil
	isErr   = convey.ShouldBeError
	isTrue  = convey.ShouldBeTrue
	isFalse = convey.ShouldBeFalse
	isZero  = convey.ShouldBeZeroValue

	hasSubStr   = convey.ShouldContainSubstring
	shouldPanic = convey.ShouldPanic
)

// go test -v -failfast -cover -coverprofile cover.out && go tool cover -html cover.out -o cover.html

func test(t *testing.T, scene string, f func(*testing.T)) {
	if t.Failed() {
		return
	}
	cv(scene, t, func() {
		f(t)
	})
}

func printBytes(t *testing.T, b []byte, prefix ...string) {
	var s string

	if len(prefix) > 0 {
		s = prefix[0]
	}
	s = s + string(b[:])
	t.Log(s)
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestJsonvalue(t *testing.T) {
	test(t, "test options", testOption)
	test(t, "test Get", testGet)
	test(t, "test Set", testSet)
	test(t, "test NewXxx", testNewXxx)
	test(t, "jsonvalue basic function", testBasicFunction)
	test(t, "misc strange characters", testMiscCharacters)
	test(t, "MustUnmarshalXxxx errors", testMustUnmarshalErrors)
	test(t, "misc simple unmarshal errors", testMiscUnmarshalErrors)
	test(t, "UTF-16 string", testUTF16)
	test(t, "percentage symbol", testPercentage)
	test(t, "misc number typed parameter", testMiscInt)
	test(t, "test an internal struct", testUnmarshalWithIter)
	test(t, "test iteration", testIteration)
	test(t, "test iteration internally", testIter)
	test(t, "test iterate float internally", testIterFloat)
	test(t, "test marshaling", testMarshal)
	test(t, "test sort", testSort)
	test(t, "test insert, append, delete", testInsertAppendDelete)
	test(t, "test import/export", testImportExport)
	test(t, "test Equal functions", testEqual)
	test(t, "test marshaler and unmarshaler", testMarshalerUnmarshaler)
	test(t, "test internal variables", testInternal)
}

func testBasicFunction(t *testing.T) {
	raw := `{"message":"hello, 世界","float":1234.123456789123456789,"true":true,"false":false,"null":null,"obj":{"msg":"hi"},"arr":["你好","world",null],"uint":1234,"int":-1234}`

	v, err := Unmarshal([]byte(raw))
	so(err, isNil)

	t.Logf("OK: %+v", v)

	b, err := v.Marshal()
	so(err, isNil)

	t.Logf("marshal: '%s'", string(b))

	// can it be unmarshal back?
	j := make(map[string]any)
	err = json.Unmarshal(b, &j)
	so(err, isNil)

	b, _ = json.Marshal(&j)
	t.Logf("marshal back: %v", string(b))

	v = NewFloat64(math.NaN())
	so(v.String(), eq, "NaN")

	v = NewFloat64(math.Inf(1))
	so(v.String(), eq, "+Inf")

	v = NewFloat64(math.Inf(-1))
	so(v.String(), eq, "-Inf")
}

func testMiscCharacters(t *testing.T) {

	cv("unmarshal and marshal back", func() {
		s := "\"/\b\f\t\r\n<>&你好世界Café\\n"
		expected := "\"\\\"\\/\\b\\f\\t\\r\\n\\u003C\\u003E\\u0026\\u4F60\\u597D\\u4E16\\u754CCaf\\u00E9\\\\n\""
		v := NewString(s)
		raw, err := v.MarshalString()
		so(err, isNil)

		printBytes(t, []byte(raw), "marshaled")
		printBytes(t, []byte(expected), "expected")
		so(raw, eq, expected)

		v, err = UnmarshalString(raw)
		so(err, isNil)

		printBytes(t, []byte(s), "Original string")
		printBytes(t, []byte(v.String()), "Got string")

		raw, err = v.MarshalString()
		so(err, isNil)
		so(raw, eq, expected)
	})

	cv("unmarshal and marshal /", func() {
		s := `"/"`
		v, err := UnmarshalString(s)
		so(err, isNil)
		so(v.IsString(), isTrue)
		so(v.String(), eq, "/")

		s = `"\/"`
		v, err = UnmarshalString(s)
		so(err, isNil)
		so(v.IsString(), isTrue)
		so(v.String(), eq, "/")
	})

	cv("unmashal UTF-8 string", func() {
		s := "你好, Café😊"
		raw := `"` + s + `"`

		printBytes(t, []byte(raw))

		v, err := UnmarshalString(raw)
		so(err, isNil)
		so(v.IsString(), isTrue)
		so(v.String(), eq, s)
	})

	cv("unmarshal illegal UTF-8 string", func() {
		s := `"😊"`
		b := []byte(s)

		printBytes(t, b, "correct bytes")
		_, err := Unmarshal(b)
		so(err, isNil)

		incompleteB := b[:2]
		printBytes(t, incompleteB, "incomplete bytes")
		_, err = Unmarshal(incompleteB)
		so(err, isErr)

		b[1] |= 0x08
		printBytes(t, b, "error bytes")
		_, err = Unmarshal(b)
		so(err, isErr)

		v, err := UnmarshalString(s)
		so(err, isNil)
		res := v.MustMarshal()
		printBytes(t, res, "correct marshaled string")
	})

	cv("unmarshal illegal escaped ASCII string", func() {
		so(ValueType(56636).String(), eq, NotExist.String())
		so(ValueType(-1).String(), eq, NotExist.String())

		v, err := UnmarshalString(`"\`)
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\u00`)
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\u0GAB`)
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\U1234`)
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\uD83d\uDE0A你好Café😊"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isNil)
		so(v.String(), eq, "😊你好Café😊")
		so(v.ValueType(), ne, NotExist)

		v, err = UnmarshalString(`"\uD83D\uDE0A, smile!"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isNil)
		so(v.ValueType(), ne, NotExist)

		v, err = UnmarshalString(`"\uD83D\uDE0"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\uD83D\UDE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\uD83D/uDE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\uD83D\uHE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\uH83D\uDE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\uD83D\u000A"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)

		v, err = UnmarshalString(`"\uD83D\uFFFF"`) // should be "\uD83D\uDE0A" ==> 😊
		so(err, isErr)
		so(v.ValueType(), eq, NotExist)
	})

	cv("unmarshal illegal plus symbols", func() {
		okCases := []string{
			`{"number":1}`,
			`{"number":1E+1}`,
		}
		failCases := []string{
			`{"number":+1}`,
			`{"number":1+1}`,
			`{"number":1E1+1}`,
			`{"number":1E1+}`,
		}

		for _, c := range okCases {
			v, err := UnmarshalString(c)
			so(err, isNil)
			so(v.MustGet("number").ValueType(), eq, Number)
		}
		for _, c := range failCases {
			_, err := UnmarshalString(c)
			so(err, isErr)
		}
	})
}

func testMustUnmarshalErrors(t *testing.T) {
	const illegal = `:\`

	v := MustUnmarshalString(illegal)
	so(v, notNil)
	so(v.ValueType(), eq, NotExist)

	v = nil
	v = MustUnmarshal([]byte(illegal))
	so(v, notNil)
	so(v.ValueType(), eq, NotExist)

	v = nil
	v = MustUnmarshalNoCopy([]byte(illegal))
	so(v, notNil)
	so(v.ValueType(), eq, NotExist)
}

func testMiscUnmarshalErrors(t *testing.T) {
	var err error

	_, err = UnmarshalString(`tru`)
	so(err, isErr)

	_, err = UnmarshalString(`fals`)
	so(err, isErr)

	_, err = UnmarshalString(`nul`)
	so(err, isErr)
}

func testUTF16(t *testing.T) {
	// orig := "你👨‍👩‍👧‍👧你"
	orig := fmt.Sprintf(
		"%c%c%c%c%c%c%c%c%c",
		0x2F804, 0x1F468, 0x200D, 0x1F469, 0x200D, 0x1F467, 0x200D, 0x1F467, 0x4F60,
	)

	v := NewObject()
	v.MustSetString(orig).At("string")

	data := struct {
		String string `json:"string"`
	}{}

	s := v.MustMarshalString()
	t.Logf("marshaled string '%s': '%s'", orig, s)
	so(orig, ne, s)

	b := v.MustMarshal()
	err := json.Unmarshal(b, &data)
	so(err, isNil)
	so(data.String, eq, orig)
}

func testPercentage(t *testing.T) {
	s := "%"
	expectedA := "\"\\u0025\""
	expectedB := "\"%\""
	v := NewString(s)
	raw, err := v.MarshalString()
	so(err, isNil)

	t.Log("marshaled: '" + raw + "'")
	so(raw != expectedA && raw != expectedB, isFalse)
}

func testMiscInt(t *testing.T) {
	var err error

	raw := `[1,2,3,4,5,6,7]`
	v, err := UnmarshalString(raw)
	so(err, isNil)

	i, err := v.GetInt(uint(2))
	so(err, isNil)
	so(i, eq, 3)

	_, err = v.GetInt(int64(2))
	so(err, isNil)

	_, err = v.GetInt(uint64(2))
	so(err, isNil)

	_, err = v.GetInt(int32(2))
	so(err, isNil)

	_, err = v.GetInt(uint32(2))
	so(err, isNil)

	_, err = v.GetInt(int16(2))
	so(err, isNil)

	_, err = v.GetInt(uint16(2))
	so(err, isNil)

	_, err = v.GetInt(int8(2))
	so(err, isNil)

	_, err = v.GetInt(uint8(2))
	so(err, isNil)
}

func testUnmarshalWithIter(t *testing.T) {
	cv("string", func() {
		raw := []byte("hello, 世界")
		rawWithQuote := []byte(fmt.Sprintf("\"%s\"", raw))

		v, err := unmarshalWithIter(globalPool{}, iter(rawWithQuote), 0)
		so(err, isNil)
		so(v.String(), eq, string(raw))
	})

	cv("true", func() {
		raw := []byte("  true  ")
		v, err := unmarshalWithIter(globalPool{}, iter(raw), 0)
		so(err, isNil)
		so(v.Bool(), isTrue)
		so(v.IsBoolean(), isTrue)
	})

	cv("false", func() {
		raw := []byte("  false  ")
		v, err := unmarshalWithIter(globalPool{}, iter(raw), 0)
		so(err, isNil)
		so(v.Bool(), isFalse)
		so(v.IsBoolean(), isTrue)
	})

	cv("null", func() {
		raw := []byte("\r\t\n  null \r\t\b  ")
		v, err := unmarshalWithIter(globalPool{}, iter(raw), 0)
		so(err, isNil)
		so(v.IsNull(), isTrue)
	})

	cv("int number", func() {
		raw := []byte(" 1234567890 ")
		v, err := unmarshalWithIter(globalPool{}, iter(raw), 0)
		so(err, isNil)
		so(v.Int64(), eq, 1234567890)
	})

	cv("array with basic type", func() {
		raw := []byte(" [123, true, false, null, [\"array in array\"], \"Hello, world!\" ] ")
		v, err := unmarshalWithIter(globalPool{}, iter(raw), 0)
		so(err, isNil)
		so(v.IsArray(), isTrue)

		t.Logf("res: %v", v)
	})

	cv("object with basic type", func() {
		raw := []byte(`  {"message": "Hello, world!"}	`)
		printBytes(t, raw)

		v, err := unmarshalWithIter(globalPool{}, iter(raw), 0)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		t.Logf("res: %v", v)

		for kv := range v.IterObjects() {
			t.Logf("key: %s", kv.K)
			t.Logf("val: %v", kv.V)
		}

		s, err := v.Get("message")
		so(err, isNil)
		so(v, notNil)
		so(s.IsString(), isTrue)
		so(s.String(), eq, "Hello, world!")
	})

	cv("object with complex type", func() {
		raw := []byte(` {"arr": [1234, true , null, false, {"obj":"empty object"}]}  `)
		printBytes(t, raw)

		v, err := unmarshalWithIter(globalPool{}, iter(raw), 0)
		so(err, isNil)
		so(v.IsObject(), isTrue)

		t.Logf("res: %v", v)

		child, err := v.Get("arr", 4, "obj")
		so(err, isNil)
		so(child.IsString(), isTrue)
		so(child.String(), eq, "empty object")
	})

	cv("nil jsonvalue object String", func() {
		var invalid *V
		s := invalid.String()
		so(s, eq, "nil")
	})
}
