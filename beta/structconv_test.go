package beta

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

// go test -v -failfast -cover -coverprofile cover.out && go tool cover -html cover.out -o cover.html

func TestJSONValueBeta(t *testing.T) {
	Convey("test structconv.go", t, func() { testStructConv(t) })
}

func testStructConv(t *testing.T) {
	Convey("test structconv.go Import", func() { testStructConv_Import(t) })
}

func testStructConv_Import(t *testing.T) {
	Convey("[]byte, json.RawMessage", func() { testStructConv_Import_RawAndBytes(t) })
	Convey("uintptr", func() { testStructConv_Import_StrangeButSupportedTypes(t) })
	Convey("invalid types", func() { testStructConv_Import_InvalidTypes(t) })
	Convey("general types", func() { testStructConv_Import_NormalTypes(t) })
	Convey("array and slice", func() { testStructConv_Import_ArrayAndSlice(t) })
}

func testStructConv_Import_RawAndBytes(t *testing.T) {
	Convey("json.RawMessage", func() {
		msg := "Hello, raw message!"
		st := struct {
			Raw json.RawMessage `json:"raw"`
		}{
			Raw: []byte(fmt.Sprintf(`{"message":"%s"}`, msg)),
		}

		j, err := Import(&st)
		So(err, ShouldBeNil)
		So(j.Len(), ShouldEqual, 1)

		got, err := j.GetString("raw", "message")
		So(err, ShouldBeNil)
		So(got, ShouldEqual, msg)

		// t.Logf("%v", j)
	})

	Convey("json.RawMessage error", func() {
		msg := "Hello, raw message!"
		st := struct {
			Raw json.RawMessage `json:"raw"`
		}{
			Raw: []byte(fmt.Sprintf(`{"message":"%s"`, msg)),
		}

		j, err := Import(&st)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)
	})

	Convey("[]byte", func() {
		b := []byte{1, 2, 3, 4}
		st := struct {
			Bytes []byte `json:"bytes"`
		}{
			Bytes: b,
		}

		j, err := Import(&st)
		So(err, ShouldBeNil)

		gotS, err := j.GetString("bytes")
		So(err, ShouldBeNil)

		gotB, err := j.GetBytes("bytes")
		So(err, ShouldBeNil)

		gotStoB, err := base64.StdEncoding.DecodeString(gotS)
		So(err, ShouldBeNil)

		So(bytes.Equal(gotB, b), ShouldBeTrue)
		So(bytes.Equal(gotStoB, b), ShouldBeTrue)

		// t.Logf("%v", j)
	})
}

func testStructConv_Import_StrangeButSupportedTypes(t *testing.T) {
	Convey("uintptr", func() {
		st := struct {
			Ptr uintptr
		}{
			Ptr: 1234,
		}

		b, err := json.Marshal(&st)
		t.Logf("Got bytes: '%s'", b)
		So(err, ShouldBeNil)

		j, err := Import(&st)
		So(err, ShouldBeNil)

		bb := j.MustMarshal()
		So(bytes.Equal(b, bb), ShouldBeTrue)
	})

	Convey("map[uintptr]xxx", func() {
		m := map[uintptr]int{
			1: 2,
			2: 3,
		}

		j, err := Import(m)
		So(err, ShouldBeNil)
		So(j.IsObject(), ShouldBeTrue)
		So(j.MustGet("1").Uint(), ShouldEqual, m[1])
		So(j.MustGet("2").Uint(), ShouldEqual, m[2])
	})

	Convey("map[int]xxx", func() {
		m := map[int]int{
			1: 2,
			2: 3,
		}

		j, err := Import(m)
		So(err, ShouldBeNil)
		So(j.IsObject(), ShouldBeTrue)
		So(j.MustGet("1").Int(), ShouldEqual, m[1])
		So(j.MustGet("2").Int(), ShouldEqual, m[2])
	})
}

func testStructConv_Import_InvalidTypes(t *testing.T) {
	Convey("complex", func() {
		st := struct {
			C complex128
		}{
			C: complex(1, 2),
		}

		_, err := json.Marshal(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)

		_, err = Import(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)
	})

	Convey("chan", func() {
		st := struct {
			Ch chan struct{}
		}{
			Ch: make(chan struct{}),
		}

		_, err := json.Marshal(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)

		_, err = Import(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)
	})

	Convey("unsafe.Pointer", func() {
		st := struct {
			Ptr unsafe.Pointer
		}{
			Ptr: nil,
		}

		_, err := json.Marshal(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)

		_, err = Import(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)
	})

	Convey("func", func() {
		st := struct {
			Func func()
		}{
			Func: func() { panic("Hey!") },
		}

		_, err := json.Marshal(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)

		_, err = Import(&st)
		So(err, ShouldBeError)
		t.Logf("expect error: %v", err)
	})

	Convey("not a struct", func() {
		ch := make(chan error)
		defer close(ch)

		j, err := Import(&ch)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)

		j, err = Import(ch)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)
	})

	Convey("map[float64]xxx", func() {
		m := map[float64]int{
			1: 2,
			2: 3,
		}

		j, err := Import(m)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)
	})

	Convey("illegal type in slice", func() {
		arr := []interface{}{
			1, complex(1, 2),
		}
		j, err := Import(arr)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)

		arr = []interface{}{
			1, []interface{}{complex(1, 2)},
		}
		j, err = Import(arr)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)
	})

	Convey("illegal type in map", func() {
		m := map[string]interface{}{
			"complex": complex(1, 2),
		}
		j, err := Import(m)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)

		m = map[string]interface{}{
			"obj": map[string]interface{}{
				"complex": complex(1, 2),
			},
		}
		j, err = Import(m)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)
	})
}

func testStructConv_Import_NormalTypes(t *testing.T) {
	Convey("bool", func() {
		st := struct {
			True   bool `json:"true"`
			False  bool `json:"false"`
			Empty  bool `json:",omitempty"`
			String bool `json:",string"`
		}{
			True:   true,
			False:  false,
			Empty:  false,
			String: true,
		}

		b, err := json.Marshal(&st)
		t.Logf("Got bytes: '%s'", b)
		So(err, ShouldBeNil)

		j, err := Import(&st)
		t.Logf("Got bytes: '%s'", j.MustMarshalString())
		So(err, ShouldBeNil)

		boo, err := j.GetBool("true")
		So(err, ShouldBeNil)
		So(boo, ShouldBeTrue)

		boo, err = j.GetBool("false")
		So(err, ShouldBeNil)
		So(boo, ShouldBeFalse)

		_, err = j.GetBool("Empty")
		So(err, ShouldBeError)

		boo, err = j.GetBool("String")
		So(err, ShouldBeError)
		So(boo, ShouldBeTrue)
	})

	Convey("number", func() {
		st := struct {
			Int   int32   `json:"int,string"`
			Uint  uint64  `json:"uint,string"`
			Float float32 `json:"float,string"`
		}{
			Int:   -100,
			Uint:  10000,
			Float: 123.125,
		}

		j, err := Import(&st)
		So(err, ShouldBeNil)

		s, err := j.GetString("int")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, strconv.Itoa(int(st.Int)))

		s, err = j.GetString("uint")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, strconv.FormatUint(uint64(st.Uint), 10))

		s, err = j.GetString("float")
		So(err, ShouldBeNil)
		So(s, ShouldEqual, strconv.FormatFloat(float64(st.Float), 'f', -1, 32))
	})
}

func testStructConv_Import_ArrayAndSlice(t *testing.T) {
	Convey("slice", func() {
		st := []struct {
			S string `json:"string"`
			I int    `json:"int"`
		}{
			{
				S: "Hello, 01",
				I: 1,
			}, {
				S: "Hello, 02",
				I: 2,
			},
		}

		j, err := Import(&st)
		So(err, ShouldBeNil)

		// t.Logf("%s", j.MustMarshalString())

		So(j.IsArray(), ShouldBeTrue)
		So(j.Len(), ShouldEqual, 2)

		for it := range j.IterArray() {
			i := it.I

			s, err := j.GetString(i, "string")
			So(err, ShouldBeNil)
			So(s, ShouldEqual, st[i].S)

			n, err := j.GetInt(i, "int")
			So(err, ShouldBeNil)
			So(n, ShouldEqual, st[i].I)
		}
	})

	Convey("array", func() {
		arr := [6]rune{'你', '好', 'J', 'S', 'O', 'N'}

		j, err := Import(&arr)
		So(err, ShouldBeNil)
		So(j.IsArray(), ShouldBeTrue)

		So(j.Len(), ShouldEqual, len(arr))
		for i, r := range arr {
			child, err := j.Get(i)
			So(err, ShouldBeNil)
			So(child.IsNumber(), ShouldBeTrue)
			So(child.Uint(), ShouldEqual, r)
		}
	})

	Convey("map[string]interface{}", func() {
		m := map[string]interface{}{
			"uint":   uint8(255),
			"float":  float32(-0.25),
			"string": "Hello, interface{}",
			"bool":   true,
			"struct": struct{}{},
		}

		j, err := Import(m)
		So(err, ShouldBeNil)
		So(j.IsObject(), ShouldBeTrue)

		// t.Log(j.MustMarshalString())

		So(j.MustGet("uint").Uint(), ShouldEqual, m["uint"])
		So(j.MustGet("float").Float64(), ShouldEqual, m["float"])
		So(j.MustGet("string").String(), ShouldEqual, m["string"])
		So(j.MustGet("bool").Bool(), ShouldEqual, m["bool"])

		st, err := j.Get("struct")
		So(err, ShouldBeNil)
		So(st.IsObject(), ShouldBeTrue)
		So(st.Len(), ShouldBeZeroValue)
	})

	Convey("anonymous struct", func() {
		type inner struct {
			F float64
		}
		st := struct {
			inner
			I int
		}{}

		st.F = 0.25
		st.I = 1024

		j, err := Import(&st)
		So(err, ShouldBeNil)

		f, err := j.GetFloat32("F")
		So(err, ShouldBeNil)
		So(f, ShouldEqual, st.F)

		i, err := j.GetFloat32("I")
		So(err, ShouldBeNil)
		So(i, ShouldEqual, st.I)
	})

	Convey("illegal anonymous struct", func() {
		type inner struct {
			Ch chan error
		}
		st := struct {
			inner
			I int
		}{}

		j, err := Import(&st)
		So(err, ShouldBeError)
		So(j, ShouldNotBeNil)
	})

	Convey("misc omitempty", func() {
		st := struct {
			I    int             `json:",omitempty"`
			U    uint            `json:",omitempty"`
			S    string          `json:",omitempty"`
			F    float64         `json:",omitempty"`
			A    []int           `json:",omitempty"`
			B    bool            `json:",omitempty"`
			By   []byte          `json:",omitempty"`
			Rw   json.RawMessage `json:",omitempty"`
			St   *struct{}       `json:",omitempty"`
			M    map[int]int     `json:",omitempty"`
			Null *int            `json:"null"`

			priv   int
			Ignore bool `json:"-"`
		}{}

		b, err := json.Marshal(&st)
		So(err, ShouldBeNil)
		s := string(b)

		j, err := Import(&st)
		So(err, ShouldBeNil)
		So(j.MustMarshalString(), ShouldEqual, s)

		st.M = map[int]int{}
		st.A = []int{}

		j, err = Import(&st)
		So(err, ShouldBeNil)
		So(j.MustMarshalString(), ShouldEqual, s)

		t.Log(s)
	})

	Convey("not omitempty", func() {
		st := struct {
			A    []int
			B    bool
			By   []byte
			St   *struct{}
			M    map[int]int
			Null *int
		}{
			A:  []int{},
			By: []byte{},
			M:  map[int]int{},
		}

		j, err := Import(&st)
		So(err, ShouldBeNil)

		So(j.MustGet("A").IsArray(), ShouldBeTrue)
		So(j.MustGet("B").IsBoolean(), ShouldBeTrue)
		So(j.MustGet("By").IsString(), ShouldBeTrue)
		So(j.MustGet("St").IsNull(), ShouldBeTrue)
		So(j.MustGet("M").IsObject(), ShouldBeTrue)
		So(j.MustGet("Null").IsNull(), ShouldBeTrue)
	})
}
