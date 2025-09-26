package jsonvalue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
)

func testMarshal(t *testing.T) {
	cv("NaN", func() { testMarshalFloat64NaN(t) })
	cv("Inf", func() { testMarshalFloat64Inf(t) })
	cv("escapeHTML", func() { testMarshalEscapeHTML(t) })
	cv("UTF-8", func() { testMarshalEscapeUTF8(t) })
	cv("slash", func() { testMarshalEscapeSlash(t) })
	cv("indent", func() { testMarshalIndent(t) })
	cv("ASCII control characters", func() { testMarshalControlCharacters(t) })
	cv("test JSONP and control ASCII for UTF-8", func() { testMarshalJSONPAndControlAsciiForUTF8(t) })
	cv("Issue #30", func() { testIssue30(t) })
	cv("MarshalWrite", func() { testMarshalWrite(t) })
}

func testMarshalFloat64NaN(*testing.T) {
	cv("with error", func() {
		v := New(math.NaN())
		_, err := v.Marshal()
		so(err, isErr)

		v = NewFloat32(float32(math.NaN()))
		_, err = v.MarshalString()
		so(err, isErr)

		v = NewFloat64(math.NaN())
		_, err = v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNHandleType(80),
		})
		so(err, isErr)
	})

	cv("to float", func() {
		v := NewFloat64(math.NaN())
		b, err := v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    1.5,
		})
		so(err, isNil)
		so(string(b), eq, "1.5")

		b, err = v.Marshal(OptFloatNaNToFloat(1.5))
		so(err, isNil)
		so(string(b), eq, "1.5")
	})

	cv("to string", func() {
		v := NewFloat64(math.NaN())
		s, err := v.MarshalString(Opt{
			FloatNaNHandleType: FloatNaNConvertToString,
		})
		so(err, isNil)
		so(s, eq, `"NaN"`)

		s, err = v.MarshalString(OptFloatNaNToStringNaN())
		so(err, isNil)
		so(s, eq, `"NaN"`)

		s, err = v.MarshalString(Opt{
			FloatNaNHandleType: FloatNaNConvertToString,
			FloatNaNToString:   "not a number",
		})
		so(err, isNil)
		so(s, eq, `"not a number"`)

		s, err = v.MarshalString(OptFloatNaNToString("not a number"))
		so(err, isNil)
		so(s, eq, `"not a number"`)
	})

	cv("to null", func() {
		v := NewFloat64(math.NaN())
		s, err := v.MarshalString(Opt{
			FloatNaNHandleType: FloatNaNNull,
		})
		so(err, isNil)
		so(s, eq, "null")

		s, err = v.MarshalString(OptFloatNaNToNull())
		so(err, isNil)
		so(s, eq, "null")
	})

	cv("to float error", func() {
		v := NewFloat64(math.NaN())
		_, err := v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    math.NaN(),
		})
		so(err, isErr)

		_, err = v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    math.Inf(1),
		})
		so(err, isErr)

		_, err = v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    math.Inf(-1),
		})
		so(err, isErr)
	})
}

func testMarshalFloat64Inf(*testing.T) {
	cv("with error", func() {
		v := NewFloat64(math.Inf(1))
		_, err := v.Marshal()
		so(err, isErr)

		v = NewFloat32(float32(math.Inf(-1)))
		_, err = v.MarshalString()
		so(err, isErr)

		v = NewFloat64(math.Inf(1))
		_, err = v.Marshal(Opt{
			FloatInfHandleType: FloatInfHandleType(80),
		})
		so(err, isErr)

		v = NewFloat64(math.Inf(-1))
		_, err = v.Marshal(Opt{
			FloatInfHandleType: FloatInfHandleType(80),
		})
		so(err, isErr)
	})

	cv("to float", func() {
		opt := Opt{
			FloatInfHandleType: FloatInfConvertToFloat,
			FloatInfToFloat:    2.25,
		}
		v := NewObject(map[string]any{
			"+inf": math.Inf(1),
			"-inf": math.Inf(-1),
		})
		s, err := v.MarshalString(opt)
		so(err, isNil)
		so(s, hasSubStr, `"+inf":2.25`)
		so(s, hasSubStr, `"-inf":-2.25`)

		s, err = v.MarshalString(OptFloatInfToFloat(2.25))
		so(err, isNil)
		so(s, hasSubStr, `"+inf":2.25`)
		so(s, hasSubStr, `"-inf":-2.25`)
	})

	cv("to string", func() {
		v := NewFloat64(math.Inf(1))
		s, err := v.MarshalString(Opt{
			FloatInfHandleType: FloatInfConvertToString,
		})
		so(err, isNil)
		so(s, eq, `"+Inf"`)

		s, err = v.MarshalString(OptFloatInfToStringInf())
		so(err, isNil)
		so(s, eq, `"+Inf"`)

		v = NewFloat64(math.Inf(-1))
		s, err = v.MarshalString(Opt{
			FloatInfHandleType: FloatInfConvertToString,
		})
		so(err, isNil)
		so(s, eq, `"-Inf"`)

		s, err = v.MarshalString(OptFloatInfToStringInf())
		so(err, isNil)
		so(s, eq, `"-Inf"`)

		v = NewObject(map[string]any{
			"+inf": math.Inf(1),
			"-inf": math.Inf(-1),
		})

		s, err = v.MarshalString(Opt{
			FloatInfHandleType:       FloatInfConvertToString,
			FloatInfPositiveToString: "infinity",
		})
		so(err, isNil)
		so(s, hasSubStr, `"+inf":"infinity"`)
		so(s, hasSubStr, `"-inf":"-infinity"`)

		s, err = v.MarshalString(OptFloatInfToString("infinity", ""))
		so(err, isNil)
		so(s, hasSubStr, `"+inf":"infinity"`)
		so(s, hasSubStr, `"-inf":"-infinity"`)

		s, err = v.MarshalString(Opt{
			FloatInfHandleType:       FloatInfConvertToString,
			FloatInfPositiveToString: "+mugen",
		})
		so(err, isNil)
		so(s, hasSubStr, `"+inf":"+mugen"`)
		so(s, hasSubStr, `"-inf":"-mugen"`)

		s, err = v.MarshalString(OptFloatInfToString("+mugen", ""))
		so(err, isNil)
		so(s, hasSubStr, `"+inf":"+mugen"`)
		so(s, hasSubStr, `"-inf":"-mugen"`)

		s, err = v.MarshalString(Opt{
			FloatInfHandleType:       FloatInfConvertToString,
			FloatInfPositiveToString: "heaven",
			FloatInfNegativeToString: "hell",
		})
		so(err, isNil)
		so(s, hasSubStr, `"+inf":"heaven"`)
		so(s, hasSubStr, `"-inf":"hell"`)

		s, err = v.MarshalString(OptFloatInfToString("heaven", "hell"))
		so(err, isNil)
		so(s, hasSubStr, `"+inf":"heaven"`)
		so(s, hasSubStr, `"-inf":"hell"`)
	})

	cv("to null", func() {
		v := NewObject(map[string]any{
			"+inf": math.Inf(1),
			"-inf": math.Inf(-1),
		})
		s, err := v.MarshalString(Opt{
			FloatInfHandleType: FloatInfNull,
		})
		so(err, isNil)
		so(s, hasSubStr, `"+inf":null`)
		so(s, hasSubStr, `"-inf":null`)

		s, err = v.MarshalString(OptFloatInfToNull())
		so(err, isNil)
		so(s, hasSubStr, `"+inf":null`)
		so(s, hasSubStr, `"-inf":null`)
	})

	cv("to float error", func() {
		iter := func(f float64) {
			v := NewFloat64(f)
			_, err := v.Marshal(Opt{
				FloatInfHandleType: FloatInfConvertToFloat,
				FloatInfToFloat:    math.NaN(),
			})
			so(err, isErr)

			_, err = v.Marshal(Opt{
				FloatInfHandleType: FloatInfConvertToFloat,
				FloatInfToFloat:    math.Inf(1),
			})
			so(err, isErr)

			_, err = v.Marshal(Opt{
				FloatInfHandleType: FloatInfConvertToFloat,
				FloatInfToFloat:    math.Inf(-1),
			})
			so(err, isErr)
		}

		iter(math.Inf(1))
		iter(math.Inf(-1))
	})
}

func testMarshalEscapeHTML(*testing.T) {
	esc := func(s string) string {
		seq := []rune{'&', '<', '>'}
		for _, r := range seq {
			s = strings.ReplaceAll(s, string(r), fmt.Sprintf("\\u00%X", r))
		}
		return s
	}

	key := "<X>&<Y>"
	value := "<12, 34> & <56, 78>"

	v := NewObject(M{
		key: value,
	})

	cv("default escape", func() {
		s := v.MustMarshalString()
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, esc(key), esc(value)))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})

	cv("escapeHTML on", func() {
		s := v.MustMarshalString(OptEscapeHTML(true))
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, esc(key), esc(value)))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})

	cv("escapeHTML off", func() {
		s := v.MustMarshalString(OptEscapeHTML(false))
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, key, value))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})
}

func testMarshalEscapeUTF8(*testing.T) {
	htmlRunes := map[rune]struct{}{
		'<': {},
		'>': {},
		'&': {},
	}

	esc := func(buf *bytes.Buffer, r rune) {
		s := fmt.Sprintf("\\u%04X", r)
		buf.WriteString(s)
	}

	key := "<一>&<二>"
	value := "<12, 34> & <56, 78>"

	v := NewObject(M{
		key: value,
	})

	cv("default escape", func() {
		str := func(s string) string {
			buf := &bytes.Buffer{}
			for _, r := range s {
				if r > 0x7F {
					esc(buf, r)
				} else if _, exist := htmlRunes[r]; exist {
					esc(buf, r)
				} else {
					buf.WriteRune(r)
				}
			}
			return buf.String()
		}

		s := v.MustMarshalString()
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})

	cv("escapeHTML on, ascii", func() {
		str := func(s string) string {
			buf := &bytes.Buffer{}
			for _, r := range s {
				if r > 0x7F {
					esc(buf, r)
				} else if _, exist := htmlRunes[r]; exist {
					esc(buf, r)
				} else {
					buf.WriteRune(r)
				}
			}
			return buf.String()
		}

		s := v.MustMarshalString(OptEscapeHTML(true))
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})

	cv("escapeHTML off, ascii", func() {
		str := func(s string) string {
			buf := &bytes.Buffer{}
			for _, r := range s {
				if r > 0x7F {
					esc(buf, r)
				} else {
					buf.WriteRune(r)
				}
			}
			return buf.String()
		}

		s := v.MustMarshalString(OptEscapeHTML(false))
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})

	cv("escapeHTML on, UTF-8 on", func() {
		str := func(s string) string {
			buf := &bytes.Buffer{}
			for _, r := range s {
				if _, exist := htmlRunes[r]; exist {
					esc(buf, r)
				} else {
					buf.WriteRune(r)
				}
			}
			return buf.String()
		}

		s := v.MustMarshalString(OptEscapeHTML(true), OptUTF8())
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})

	cv("escapeHTML off, UTF-8 on", func() {
		s := v.MustMarshalString(OptEscapeHTML(false), OptUTF8())
		so(s, eq, fmt.Sprintf(`{"%s":"%s"}`, key, value))

		vv, err := UnmarshalString(s)
		so(err, isNil)
		so(vv.MustGet(key).String(), eq, value)
	})
}

func testMarshalEscapeSlash(*testing.T) {
	v := NewString("https://google.com")
	defaultStr := `"https:\/\/google.com"`
	nonEsc := `"https://google.com"`

	cv("default", func() {
		s := v.MustMarshalString()
		so(s, eq, defaultStr)
		so(MustUnmarshalString(s).String(), eq, v.String())
	})

	cv("escape slash", func() {
		s := v.MustMarshalString(OptEscapeSlash(true))
		so(s, eq, defaultStr)
		so(MustUnmarshalString(s).String(), eq, v.String())
	})

	cv("non-escape slash", func() {
		s := v.MustMarshalString(OptEscapeSlash(false))
		so(s, eq, nonEsc)
		so(MustUnmarshalString(s).String(), eq, v.String())
	})
}

func testMarshalIndent(t *testing.T) {
	cv("object", func() {
		v := NewObject()
		v.MustSetString("Hello, world").At("obj", "obj_in_obj", "msg")
		b := v.MustMarshal(OptIndent("", "  "))

		var m any
		_ = json.Unmarshal(b, &m)
		bJS, _ := json.MarshalIndent(m, "", "  ")

		so(string(b), eq, string(bJS))
		t.Log(string(b))

		b = v.MustMarshal(OptIndent("+", "  "))
		bJS, _ = json.MarshalIndent(m, "+", "  ")
		so(string(b), eq, string(bJS))
		t.Log(string(b))
	})

	cv("array", func() {
		v := NewArray()
		v.MustAppend(1).InTheEnd()
		v.MustAppend(2).InTheEnd()
		v.MustAppend(3).InTheEnd()

		b := v.MustMarshal(OptIndent("", "  "))

		var m any
		_ = json.Unmarshal(b, &m)
		bJS, _ := json.MarshalIndent(m, "", "  ")

		so(string(b), eq, string(bJS))
		t.Log(string(b))

		b = v.MustMarshal(OptIndent("+", "  "))
		bJS, _ = json.MarshalIndent(m, "+", "  ")
		so(string(b), eq, string(bJS))
		t.Log(string(b))
	})

	cv("multiple indents", func() {
		type s struct {
			Arr []any  `json:"arr,omitempty"`
			Obj *s     `json:"obj,omitempty"`
			Str string `json:"str,omitempty"`
		}

		data := &s{
			Str: "Lv.0",
			Obj: &s{
				Str: "Lv.1",
				Obj: &s{
					Str: "Lv.2",
				},
				Arr: []any{
					1,
					"2",
					&s{
						Str: "Lv1.1",
					},
				},
			},
		}

		v, err := Import(data)
		so(err, isNil)

		b := v.MustMarshal(OptIndent("", "  "), OptDefaultStringSequence())
		// b := v.MustMarshal(OptIndent("", "  "))
		bJS, _ := json.MarshalIndent(data, "", "  ")

		so(string(b), eq, string(bJS))
		t.Log(string(b))
	})

	cv("empty indent", func() {
		v := NewObject()
		v.MustSetString("Hello, world").At("obj", "obj_in_obj", "msg")
		b := v.MustMarshal(OptIndent("", ""))

		var m any
		_ = json.Unmarshal(b, &m)
		bJS, _ := json.MarshalIndent(m, "", "")

		so(string(b), eq, string(bJS))
		t.Log(string(b))
	})
}

func testMarshalControlCharacters(t *testing.T) {
	unShowAbleControlChars := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x7F,
	}
	expectedBuilder := strings.Builder{}
	expectedBuilder.WriteByte('"')
	for _, b := range unShowAbleControlChars {
		expectedBuilder.WriteString(fmt.Sprintf("\\u%04X", b))
	}
	expectedBuilder.WriteByte('"')
	expected := expectedBuilder.String()

	goVer, _ := json.Marshal(string(unShowAbleControlChars))
	v := New(string(unShowAbleControlChars))
	s := v.MustMarshalString()
	so(s, eq, expected)
	so(strings.ToLower(s), eq, strings.ReplaceAll(string(goVer), string('\u007f'), "\\u007f"))
	// It is strange that encoding/json does not escape DEL symbol

	t.Log("")
	t.Logf("jsonvalue   result: %s", s)
	t.Logf("encoding/go result: %s", goVer)

	v, err := UnmarshalString(s)
	so(err, isNil)
	so(v.String(), eq, string(unShowAbleControlChars))
}

func testMarshalJSONPAndControlAsciiForUTF8(t *testing.T) {
	unShowAbleControlCharsAndJSONPSpecial := []rune{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x0E, 0x0F,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F,
		0x7F,
		0x2028, 0x2029,
	}

	expectedBuilder := strings.Builder{}
	expectedBuilder.WriteByte('"')
	for _, r := range unShowAbleControlCharsAndJSONPSpecial {
		expectedBuilder.WriteString(fmt.Sprintf("\\u%04X", r))
	}
	expectedBuilder.WriteByte('"')
	expected := expectedBuilder.String()

	goVer, _ := json.Marshal(string(unShowAbleControlCharsAndJSONPSpecial))
	v := New(string(unShowAbleControlCharsAndJSONPSpecial))
	s := v.MustMarshalString(OptUTF8())
	so(s, eq, expected)
	so(strings.ToLower(s), eq, strings.ReplaceAll(string(goVer), string('\u007f'), "\\u007f"))

	t.Log("")
	t.Logf("jsonvalue   result: %s", s)
	t.Logf("encoding/go result: %s", goVer)

	v, err := UnmarshalString(s)
	so(err, isNil)
	so(v.String(), eq, string(unShowAbleControlCharsAndJSONPSpecial))
}

func testIssue30(t *testing.T) {
	cv("Issue #30", func() {
		responseData, err := MustUnmarshalString(issue30Raw).Get("data")
		so(err, isNil)

		var candidateInfoArr []*V
		candidateInfoArr = append(candidateInfoArr, responseData.ForRangeArr()...)
		data := candidateInfoArr

		v := NewObject()
		v.MustSet(data).At("data")

		b := v.MustMarshal(OptSetSequence())
		t.Log(v)
		t.Log(string(b))

		var m json.RawMessage
		err = json.Unmarshal(b, &m)
		so(err, isNil)
	})

	cv("stripped Issue #30", func() {
		list := []*V{
			NewObject(map[string]any{"archived": false}),
			NewObject(map[string]any{"archived": false}),
		}

		v, err := Import(list)
		so(err, isNil)
		so(v.IsArray(), isTrue)

		t.Log(v)
	})
}

func testMarshalWrite(t *testing.T) {
	cv("basic functionality", func() { testMarshalWriteBasic(t) })
	cv("nil writer error", func() { testMarshalWriteNilWriter(t) })
	cv("consistency with Marshal", func() { testMarshalWriteConsistency(t) })
	cv("different data types", func() { testMarshalWriteDataTypes(t) })
	cv("options compatibility", func() { testMarshalWriteOptions(t) })
	cv("uninitialized value", func() { testMarshalWriteUninitialized(t) })
}

func testMarshalWriteBasic(*testing.T) {
	cv("simple string", func() {
		v := NewString("Hello, world!")
		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isNil)
		so(buf.String(), eq, `"Hello, world!"`)
	})

	cv("simple number", func() {
		v := NewFloat64(123.456)
		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isNil)
		so(buf.String(), eq, "123.456")
	})

	cv("boolean true", func() {
		v := NewBool(true)
		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isNil)
		so(buf.String(), eq, "true")
	})

	cv("boolean false", func() {
		v := NewBool(false)
		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isNil)
		so(buf.String(), eq, "false")
	})

	cv("null value", func() {
		v := NewNull()
		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isNil)
		so(buf.String(), eq, "null")
	})
}

func testMarshalWriteNilWriter(*testing.T) {
	cv("nil writer should return error", func() {
		v := NewString("test")
		err := v.MarshalWrite(nil)
		so(err, isErr)
		so(err, eq, ErrNilParameter)
	})
}

func testMarshalWriteConsistency(*testing.T) {
	testData := []struct {
		name string
		v    *V
	}{
		{"string", NewString("test string with 中文 and special chars: <>&\"'")},
		{"integer", NewInt64(42)},
		{"float", NewFloat64(3.14159)},
		{"boolean true", NewBool(true)},
		{"boolean false", NewBool(false)},
		{"null", NewNull()},
		{"empty object", NewObject()},
		{"empty array", NewArray()},
	}

	for _, td := range testData {
		cv(td.name, func() {
			// Test without options
			marshalResult, err := td.v.Marshal()
			so(err, isNil)

			buf := &bytes.Buffer{}
			err = td.v.MarshalWrite(buf)
			so(err, isNil)
			so(buf.String(), eq, string(marshalResult))

			// Test with some options
			marshalResult, err = td.v.Marshal(OptEscapeHTML(false), OptUTF8())
			so(err, isNil)

			buf.Reset()
			err = td.v.MarshalWrite(buf, OptEscapeHTML(false), OptUTF8())
			so(err, isNil)
			so(buf.String(), eq, string(marshalResult))
		})
	}
}

func testMarshalWriteDataTypes(*testing.T) {
	cv("complex object", func() {
		v := NewObject()
		v.MustSetString("John").At("name")
		v.MustSetInt(30).At("age")
		v.MustSetBool(true).At("active")
		v.MustSetNull().At("data")

		arr := NewArray()
		arr.MustAppend("item1").InTheEnd()
		arr.MustAppend(123).InTheEnd()
		v.MustSet(arr).At("items")

		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isNil)

		result := buf.String()
		so(result, hasSubStr, `"name":"John"`)
		so(result, hasSubStr, `"age":30`)
		so(result, hasSubStr, `"active":true`)
		so(result, hasSubStr, `"data":null`)
		so(result, hasSubStr, `"items":["item1",123]`)
	})

	cv("nested array", func() {
		v := NewArray()
		v.MustAppend(1).InTheEnd()

		subArr := NewArray()
		subArr.MustAppend("nested").InTheEnd()
		subArr.MustAppend(true).InTheEnd()
		v.MustAppend(subArr).InTheEnd()

		v.MustAppend(3.14).InTheEnd()

		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isNil)
		so(buf.String(), eq, `[1,["nested",true],3.14]`)
	})
}

func testMarshalWriteOptions(*testing.T) {
	cv("indent option", func() {
		v := NewObject()
		v.MustSetString("value").At("key")
		v.MustSetInt(42).At("number")

		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf, OptIndent("", "  "), OptSetSequence())
		so(err, isNil)

		expected := "{\n  \"key\": \"value\",\n  \"number\": 42\n}"
		so(buf.String(), eq, expected)
	})

	cv("escape HTML option", func() {
		v := NewString("<script>alert('test')</script>")

		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf, OptEscapeHTML(true))
		so(err, isNil)
		so(buf.String(), hasSubStr, "\\u003C")
		so(buf.String(), hasSubStr, "\\u003E")

		buf.Reset()
		err = v.MarshalWrite(buf, OptEscapeHTML(false))
		so(err, isNil)
		so(buf.String(), eq, `"<script>alert('test')<\/script>"`)
	})

	cv("UTF-8 option", func() {
		v := NewString("Hello 世界")

		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf, OptUTF8())
		so(err, isNil)
		so(buf.String(), eq, `"Hello 世界"`)

		buf.Reset()
		err = v.MarshalWrite(buf)
		so(err, isNil)
		so(buf.String(), hasSubStr, "\\u")
	})

	cv("float NaN handling", func() {
		v := NewFloat64(math.NaN())

		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf, OptFloatNaNToNull())
		so(err, isNil)
		so(buf.String(), eq, "null")

		buf.Reset()
		err = v.MarshalWrite(buf, OptFloatNaNToString("not_a_number"))
		so(err, isNil)
		so(buf.String(), eq, `"not_a_number"`)

		buf.Reset()
		err = v.MarshalWrite(buf, OptFloatNaNToFloat(0.0))
		so(err, isNil)
		so(buf.String(), eq, "0")
	})

	cv("float Inf handling", func() {
		v := NewFloat64(math.Inf(1))

		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf, OptFloatInfToNull())
		so(err, isNil)
		so(buf.String(), eq, "null")

		buf.Reset()
		err = v.MarshalWrite(buf, OptFloatInfToString("infinity", "negative_infinity"))
		so(err, isNil)
		so(buf.String(), eq, `"infinity"`)

		v = NewFloat64(math.Inf(-1))
		buf.Reset()
		err = v.MarshalWrite(buf, OptFloatInfToString("infinity", "negative_infinity"))
		so(err, isNil)
		so(buf.String(), eq, `"negative_infinity"`)
	})
}

func testMarshalWriteUninitialized(*testing.T) {
	cv("uninitialized value should return error", func() {
		v := &V{} // Create an uninitialized V (valueType will be NotExist)
		buf := &bytes.Buffer{}
		err := v.MarshalWrite(buf)
		so(err, isErr)
		so(err, eq, ErrValueUninitialized)
	})
}

const issue30Raw = `{
	"data": [
		{
			"basicInfo": {
				"num": 0,
				"str": "",
				"archived": true
			}
		},
		{
			"basicInfo": {
				"num": 1,
				"str": "",
				"archived": false
			}
		}
	]
}`
