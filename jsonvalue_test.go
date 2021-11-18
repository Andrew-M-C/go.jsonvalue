package jsonvalue

import (
	"encoding/json"
	"fmt"
	"testing"

	amcbytes "github.com/Andrew-M-C/go.util/bytes"
	. "github.com/smartystreets/goconvey/convey"
)

// go test -v -failfast -cover -coverprofile cover.out && go tool cover -html cover.out -o cover.html

func test(t *testing.T, scene string, f func(*testing.T)) {
	if t.Failed() {
		return
	}
	Convey(scene, t, func() {
		f(t)
	})
}

func printBytes(t *testing.T, b []byte, prefix ...string) {
	s := amcbytes.SDump(b, prefix...)
	t.Log(s)
}

func TestJsonvalue(t *testing.T) {
	test(t, "jsonvalue basic function", testBasicFunction)
	test(t, "misc strange characters", testMiscCharacters)
	test(t, "MustUnmarshalXxxx errors", testMustUnmarshalErrors)
	test(t, "misc simple unmarshal errors", testMiscUnmarshalErrors)
	test(t, "UTF-16 string", testUTF16)
	test(t, "percentage symbol", testPercentage)
	test(t, "misc number typed parameter", testMiscInt)
	test(t, "test an internal struct", testUnmarshalWithIter)
}

func testBasicFunction(t *testing.T) {
	raw := `{"message":"hello, 世界","float":1234.123456789123456789,"true":true,"false":false,"null":null,"obj":{"msg":"hi"},"arr":["你好","world",null],"uint":1234,"int":-1234}`

	v, err := Unmarshal([]byte(raw))
	So(err, ShouldBeNil)

	t.Logf("OK: %+v", v)

	b, err := v.Marshal()
	So(err, ShouldBeNil)

	t.Logf("marshal: '%s'", string(b))

	// can it be unmarshal back?
	j := make(map[string]interface{})
	err = json.Unmarshal(b, &j)
	So(err, ShouldBeNil)

	b, _ = json.Marshal(&j)
	t.Logf("marshal back: %v", string(b))
}

func testMiscCharacters(t *testing.T) {

	Convey("unmarshal and marshal back", func() {
		s := "\"/\b\f\t\r\n<>&你好世界Café\\n"
		expected := "\"\\\"\\/\\b\\f\\t\\r\\n\\u003C\\u003E\\u0026\\u4F60\\u597D\\u4E16\\u754CCaf\\u00E9\\\\n\""
		v := NewString(s)
		raw, err := v.MarshalString()
		So(err, ShouldBeNil)

		printBytes(t, []byte(raw), "marshaled")
		printBytes(t, []byte(expected), "expected")
		So(raw, ShouldEqual, expected)

		v, err = UnmarshalString(raw)
		So(err, ShouldBeNil)

		printBytes(t, []byte(s), "Original string")
		printBytes(t, []byte(v.String()), "Got string")

		raw, err = v.MarshalString()
		So(err, ShouldBeNil)
		So(raw, ShouldEqual, expected)
	})

	Convey("unmarshal and marshal /", func() {
		s := `"/"`
		v, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(v.IsString(), ShouldBeTrue)
		So(v.String(), ShouldEqual, "/")

		s = `"\/"`
		v, err = UnmarshalString(s)
		So(err, ShouldBeNil)
		So(v.IsString(), ShouldBeTrue)
		So(v.String(), ShouldEqual, "/")
	})

	Convey("unmashal UTF-8 string", func() {
		s := "你好, Café😊"
		raw := `"` + s + `"`

		printBytes(t, []byte(raw))

		v, err := UnmarshalString(raw)
		So(err, ShouldBeNil)
		So(v.IsString(), ShouldBeTrue)
		So(v.String(), ShouldEqual, s)
	})

	Convey("unmarshal illegal UTF-8 string", func() {
		s := `"😊"`
		b := []byte(s)

		printBytes(t, b, "correct bytes")
		_, err := Unmarshal(b)
		So(err, ShouldBeNil)

		incompleteB := b[:2]
		printBytes(t, incompleteB, "incomplete bytes")
		_, err = Unmarshal(incompleteB)
		So(err, ShouldBeError)

		b[1] |= 0x08
		printBytes(t, b, "error bytes")
		_, err = Unmarshal(b)
		So(err, ShouldBeError)

		v, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		res := v.MustMarshal()
		printBytes(t, res, "correct marshaled string")
	})

	Convey("unmarshal illegal escaped ASCII string", func() {
		So(ValueType(56636).String(), ShouldEqual, NotExist.String())
		So(ValueType(-1).String(), ShouldEqual, NotExist.String())

		v, err := UnmarshalString(`"\`)
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\u00`)
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\u0GAB`)
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\U1234`)
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\uD83d\uDE0A你好Café😊"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeNil)
		So(v.String(), ShouldEqual, "😊你好Café😊")
		So(v.ValueType(), ShouldNotEqual, NotExist)

		v, err = UnmarshalString(`"\uD83D\uDE0A, smile!"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeNil)
		So(v.ValueType(), ShouldNotEqual, NotExist)

		v, err = UnmarshalString(`"\uD83D\uDE0"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\uD83D\UDE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\uD83D/uDE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\uD83D\uHE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\uH83D\uDE0A"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\uD83D\u000A"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)

		v, err = UnmarshalString(`"\uD83D\uFFFF"`) // should be "\uD83D\uDE0A" ==> 😊
		So(err, ShouldBeError)
		So(v.ValueType(), ShouldEqual, NotExist)
	})
}

func testMustUnmarshalErrors(t *testing.T) {
	const illegal = `:\`

	v := MustUnmarshalString(illegal)
	So(v, ShouldNotBeNil)
	So(v.ValueType(), ShouldEqual, NotExist)

	v = nil
	v = MustUnmarshal([]byte(illegal))
	So(v, ShouldNotBeNil)
	So(v.ValueType(), ShouldEqual, NotExist)

	v = nil
	v = MustUnmarshalNoCopy([]byte(illegal))
	So(v, ShouldNotBeNil)
	So(v.ValueType(), ShouldEqual, NotExist)
}

func testMiscUnmarshalErrors(t *testing.T) {
	var err error

	_, err = UnmarshalString(`tru`)
	So(err, ShouldBeError)

	_, err = UnmarshalString(`fals`)
	So(err, ShouldBeError)

	_, err = UnmarshalString(`nul`)
	So(err, ShouldBeError)
}

func testUTF16(t *testing.T) {
	// orig := "你👨‍👩‍👧‍👧你"
	orig := fmt.Sprintf(
		"%c%c%c%c%c%c%c%c%c",
		0x2F804, 0x1F468, 0x200D, 0x1F469, 0x200D, 0x1F467, 0x200D, 0x1F467, 0x4F60,
	)

	v := NewObject()
	v.SetString(orig).At("string")

	data := struct {
		String string `json:"string"`
	}{}

	s := v.MustMarshalString()
	t.Logf("marshaled string '%s': '%s'", orig, s)
	So(orig, ShouldNotEqual, s)

	b := v.MustMarshal()
	err := json.Unmarshal(b, &data)
	So(err, ShouldBeNil)
	So(data.String, ShouldEqual, orig)
}

func testPercentage(t *testing.T) {
	s := "%"
	expectedA := "\"\\u0025\""
	expectedB := "\"%\""
	v := NewString(s)
	raw, err := v.MarshalString()
	So(err, ShouldBeNil)

	t.Log("marshaled: '" + raw + "'")
	So(raw != expectedA && raw != expectedB, ShouldBeFalse)
}

func testMiscInt(t *testing.T) {
	var err error

	raw := `[1,2,3,4,5,6,7]`
	v, err := UnmarshalString(raw)
	So(err, ShouldBeNil)

	i, err := v.GetInt(uint(2))
	So(err, ShouldBeNil)
	So(i, ShouldEqual, 3)

	_, err = v.GetInt(int64(2))
	So(err, ShouldBeNil)

	_, err = v.GetInt(uint64(2))
	So(err, ShouldBeNil)

	_, err = v.GetInt(int32(2))
	So(err, ShouldBeNil)

	_, err = v.GetInt(uint32(2))
	So(err, ShouldBeNil)

	_, err = v.GetInt(int16(2))
	So(err, ShouldBeNil)

	_, err = v.GetInt(uint16(2))
	So(err, ShouldBeNil)

	_, err = v.GetInt(int8(2))
	So(err, ShouldBeNil)

	_, err = v.GetInt(uint8(2))
	So(err, ShouldBeNil)
}

func testUnmarshalWithIter(t *testing.T) {
	Convey("string", func() {
		raw := []byte("hello, 世界")
		rawWithQuote := []byte(fmt.Sprintf("\"%s\"", raw))

		v, err := unmarshalWithIter(&iter{b: rawWithQuote}, 0)
		So(err, ShouldBeNil)
		So(v.String(), ShouldEqual, string(raw))
	})

	Convey("true", func() {
		raw := []byte("  true  ")
		v, err := unmarshalWithIter(&iter{b: raw}, 0)
		So(err, ShouldBeNil)
		So(v.Bool(), ShouldBeTrue)
		So(v.IsBoolean(), ShouldBeTrue)
	})

	Convey("false", func() {
		raw := []byte("  false  ")
		v, err := unmarshalWithIter(&iter{b: raw}, 0)
		So(err, ShouldBeNil)
		So(v.Bool(), ShouldBeFalse)
		So(v.IsBoolean(), ShouldBeTrue)
	})

	Convey("null", func() {
		raw := []byte("\r\t\n  null \r\t\b  ")
		v, err := unmarshalWithIter(&iter{b: raw}, 0)
		So(err, ShouldBeNil)
		So(v.IsNull(), ShouldBeTrue)
	})

	Convey("int number", func() {
		raw := []byte(" 1234567890 ")
		v, err := unmarshalWithIter(&iter{b: raw}, 0)
		So(err, ShouldBeNil)
		So(v.Int64(), ShouldEqual, 1234567890)
	})

	Convey("array with basic type", func() {
		raw := []byte(" [123, true, false, null, [\"array in array\"], \"Hello, world!\" ] ")
		v, err := unmarshalWithIter(&iter{b: raw}, 0)
		So(err, ShouldBeNil)
		So(v.IsArray(), ShouldBeTrue)

		t.Logf("res: %v", v)
	})

	Convey("object with basic type", func() {
		raw := []byte(`  {"message": "Hello, world!"}	`)
		printBytes(t, raw)

		v, err := unmarshalWithIter(&iter{b: raw}, 0)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		t.Logf("res: %v", v)

		for kv := range v.IterObjects() {
			t.Logf("key: %s", kv.K)
			t.Logf("val: %v", kv.V)
		}

		s, err := v.Get("message")
		So(err, ShouldBeNil)
		So(v, ShouldNotBeNil)
		So(s.IsString(), ShouldBeTrue)
		So(s.String(), ShouldEqual, "Hello, world!")
	})

	Convey("object with complex type", func() {
		raw := []byte(` {"arr": [1234, true , null, false, {"obj":"empty object"}]}  `)
		printBytes(t, raw)

		v, err := unmarshalWithIter(&iter{b: raw}, 0)
		So(err, ShouldBeNil)
		So(v.IsObject(), ShouldBeTrue)

		t.Logf("res: %v", v)

		child, err := v.Get("arr", 4, "obj")
		So(err, ShouldBeNil)
		So(child.IsString(), ShouldBeTrue)
		So(child.String(), ShouldEqual, "empty object")
	})
}
