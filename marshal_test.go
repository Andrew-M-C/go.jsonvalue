package jsonvalue

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshal(t *testing.T) {
	test(t, "NaN", testMarshalFloat64NaN)
	test(t, "Inf", testMarshalFloat64Inf)
	test(t, "escapeHTML", testMarshalEscapeHTML)
	test(t, "UTF-8", testMarshalEscapeUTF8)
}

func testMarshalFloat64NaN(t *testing.T) {
	Convey("with error", func() {
		v := NewFloat64(math.NaN())
		_, err := v.Marshal()
		So(err, ShouldBeError)

		v = NewFloat32(float32(math.NaN()))
		_, err = v.MarshalString()
		So(err, ShouldBeError)

		v = NewFloat64(math.NaN())
		_, err = v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNHandleType(80),
		})
		So(err, ShouldBeError)
	})

	Convey("to float", func() {
		v := NewFloat64(math.NaN())
		b, err := v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    1.5,
		})
		So(err, ShouldBeNil)
		So(string(b), ShouldEqual, "1.5")

		b, err = v.Marshal(OptFloatNaNToFloat(1.5))
		So(err, ShouldBeNil)
		So(string(b), ShouldEqual, "1.5")
	})

	Convey("to string", func() {
		v := NewFloat64(math.NaN())
		s, err := v.MarshalString(Opt{
			FloatNaNHandleType: FloatNaNConvertToString,
		})
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"NaN"`)

		s, err = v.MarshalString(OptFloatNaNToStringNaN())
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"NaN"`)

		s, err = v.MarshalString(Opt{
			FloatNaNHandleType: FloatNaNConvertToString,
			FloatNaNToString:   "not a number",
		})
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"not a number"`)

		s, err = v.MarshalString(OptFloatNaNToString("not a number"))
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"not a number"`)
	})

	Convey("to null", func() {
		v := NewFloat64(math.NaN())
		s, err := v.MarshalString(Opt{
			FloatNaNHandleType: FloatNaNNull,
		})
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "null")

		s, err = v.MarshalString(OptFloatNaNToNull())
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "null")
	})

	Convey("to float error", func() {
		v := NewFloat64(math.NaN())
		_, err := v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    math.NaN(),
		})
		So(err, ShouldBeError)

		_, err = v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    math.Inf(1),
		})
		So(err, ShouldBeError)

		_, err = v.Marshal(Opt{
			FloatNaNHandleType: FloatNaNConvertToFloat,
			FloatNaNToFloat:    math.Inf(-1),
		})
		So(err, ShouldBeError)
	})
}

func testMarshalFloat64Inf(t *testing.T) {
	Convey("with error", func() {
		v := NewFloat64(math.Inf(1))
		_, err := v.Marshal()
		So(err, ShouldBeError)

		v = NewFloat32(float32(math.Inf(-1)))
		_, err = v.MarshalString()
		So(err, ShouldBeError)

		v = NewFloat64(math.Inf(1))
		_, err = v.Marshal(Opt{
			FloatInfHandleType: FloatInfHandleType(80),
		})
		So(err, ShouldBeError)

		v = NewFloat64(math.Inf(-1))
		_, err = v.Marshal(Opt{
			FloatInfHandleType: FloatInfHandleType(80),
		})
		So(err, ShouldBeError)
	})

	Convey("to float", func() {
		opt := Opt{
			FloatInfHandleType: FloatInfConvertToFloat,
			FloatInfToFloat:    2.25,
		}
		v := NewObject(map[string]interface{}{
			"+inf": math.Inf(1),
			"-inf": math.Inf(-1),
		})
		s, err := v.MarshalString(opt)
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":2.25`)
		So(s, ShouldContainSubstring, `"-inf":-2.25`)

		s, err = v.MarshalString(OptFloatInfToFloat(2.25))
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":2.25`)
		So(s, ShouldContainSubstring, `"-inf":-2.25`)
	})

	Convey("to string", func() {
		v := NewFloat64(math.Inf(1))
		s, err := v.MarshalString(Opt{
			FloatInfHandleType: FloatInfConvertToString,
		})
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"+Inf"`)

		s, err = v.MarshalString(OptFloatInfToStringInf())
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"+Inf"`)

		v = NewFloat64(math.Inf(-1))
		s, err = v.MarshalString(Opt{
			FloatInfHandleType: FloatInfConvertToString,
		})
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"-Inf"`)

		s, err = v.MarshalString(OptFloatInfToStringInf())
		So(err, ShouldBeNil)
		So(s, ShouldEqual, `"-Inf"`)

		v = NewObject(map[string]interface{}{
			"+inf": math.Inf(1),
			"-inf": math.Inf(-1),
		})

		s, err = v.MarshalString(Opt{
			FloatInfHandleType:       FloatInfConvertToString,
			FloatInfPositiveToString: "infinity",
		})
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":"infinity"`)
		So(s, ShouldContainSubstring, `"-inf":"-infinity"`)

		s, err = v.MarshalString(OptFloatInfToString("infinity", ""))
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":"infinity"`)
		So(s, ShouldContainSubstring, `"-inf":"-infinity"`)

		s, err = v.MarshalString(Opt{
			FloatInfHandleType:       FloatInfConvertToString,
			FloatInfPositiveToString: "+mugen",
		})
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":"+mugen"`)
		So(s, ShouldContainSubstring, `"-inf":"-mugen"`)

		s, err = v.MarshalString(OptFloatInfToString("+mugen", ""))
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":"+mugen"`)
		So(s, ShouldContainSubstring, `"-inf":"-mugen"`)

		s, err = v.MarshalString(Opt{
			FloatInfHandleType:       FloatInfConvertToString,
			FloatInfPositiveToString: "heaven",
			FloatInfNegativeToString: "hell",
		})
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":"heaven"`)
		So(s, ShouldContainSubstring, `"-inf":"hell"`)

		s, err = v.MarshalString(OptFloatInfToString("heaven", "hell"))
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":"heaven"`)
		So(s, ShouldContainSubstring, `"-inf":"hell"`)
	})

	Convey("to null", func() {
		v := NewObject(map[string]interface{}{
			"+inf": math.Inf(1),
			"-inf": math.Inf(-1),
		})
		s, err := v.MarshalString(Opt{
			FloatInfHandleType: FloatInfNull,
		})
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":null`)
		So(s, ShouldContainSubstring, `"-inf":null`)

		s, err = v.MarshalString(OptFloatInfToNull())
		So(err, ShouldBeNil)
		So(s, ShouldContainSubstring, `"+inf":null`)
		So(s, ShouldContainSubstring, `"-inf":null`)
	})

	Convey("to float error", func() {
		iter := func(f float64) {
			v := NewFloat64(f)
			_, err := v.Marshal(Opt{
				FloatInfHandleType: FloatInfConvertToFloat,
				FloatInfToFloat:    math.NaN(),
			})
			So(err, ShouldBeError)

			_, err = v.Marshal(Opt{
				FloatInfHandleType: FloatInfConvertToFloat,
				FloatInfToFloat:    math.Inf(1),
			})
			So(err, ShouldBeError)

			_, err = v.Marshal(Opt{
				FloatInfHandleType: FloatInfConvertToFloat,
				FloatInfToFloat:    math.Inf(-1),
			})
			So(err, ShouldBeError)
		}

		iter(math.Inf(1))
		iter(math.Inf(-1))
	})
}

func testMarshalEscapeHTML(t *testing.T) {
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

	Convey("default escape", func() {
		s := v.MustMarshalString()
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, esc(key), esc(value)))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})

	Convey("escapeHTML on", func() {
		s := v.MustMarshalString(OptEscapeHTML(true))
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, esc(key), esc(value)))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})

	Convey("escapeHTML off", func() {
		s := v.MustMarshalString(OptEscapeHTML(false))
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, key, value))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})
}

func testMarshalEscapeUTF8(t *testing.T) {
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

	Convey("default escape", func() {
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
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})

	Convey("escapeHTML on, ascii", func() {
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
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})

	Convey("escapeHTML off, ascii", func() {
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
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})

	Convey("escapeHTML on, UTF-8 on", func() {
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
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, str(key), str(value)))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})

	Convey("escapeHTML off, UTF-8 on", func() {
		s := v.MustMarshalString(OptEscapeHTML(false), OptUTF8())
		So(s, ShouldEqual, fmt.Sprintf(`{"%s":"%s"}`, key, value))

		vv, err := UnmarshalString(s)
		So(err, ShouldBeNil)
		So(vv.MustGet(key).String(), ShouldEqual, value)
	})
}
