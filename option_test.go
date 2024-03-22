package jsonvalue

import (
	"encoding/json"
	"strconv"
	"testing"
)

// This function used to ovwrwrite system default marshal options
func init() {
	SetDefaultMarshalOptions(OptEscapeSlash(false))
}

func testOption(t *testing.T) {
	cv("test default option overwriting", func() { testOptionOverwriting(t) })
	cv("test reset marshal options", func() { testOptionReset(t) })
	cv("test OptSetSequence", func() { testOption_OptSetSequence(t) })
	cv("test OptIgnoreOmitempty", func() { testOption_OptIgnoreOmitempty(t) })
	cv("test Issue #29", func() { testOption_Issue29(t) })
}

func testOptionOverwriting(*testing.T) {
	v := NewObject(M{
		"slash": "/",
	})

	expect := `{"slash":"/"}`
	s := v.MustMarshalString()
	so(s, eq, expect)
}

func testOptionReset(*testing.T) {
	raw := `{"slash":"/"}`
	esc := `{"slash":"\/"}`

	v := MustUnmarshalString(raw)
	so(v.IsObject(), isTrue)

	s := v.MustMarshalString()
	so(s, eq, raw)

	ResetDefaultMarshalOptions()
	s = v.MustMarshalString()
	so(s, eq, esc)

	SetDefaultMarshalOptions()
	s = v.MustMarshalString()
	so(s, eq, esc)
}

func testOption_OptSetSequence(*testing.T) {
	cv("by set", func() {
		v := NewObject()
		const total = 10
		const iterate = 1000

		for i := 0; i < total; i++ {
			v.MustSet(i).At(strconv.Itoa(i))
		}

		so(v.Len(), ne, 0)
		so(v.Len(), eq, total)

		expected := v.MustMarshalString(OptSetSequence())
		so(expected, ne, "")
		so(expected, ne, "{}")

		for i := 0; i < iterate; i++ {
			s := v.MustMarshalString(OptSetSequence())
			so(s, eq, expected)
		}
	})

	cv("by unmarshal", func() {
		const iterate = 1000
		const raw = `{"one":1,"two":2,"three":3}`
		v := MustUnmarshalString(raw)

		so(v.Len(), eq, 3)

		for i := 0; i < iterate; i++ {
			s := v.MustMarshalString(OptSetSequence())
			so(s, eq, raw)
		}
	})
}

func testOption_OptIgnoreOmitempty(t *testing.T) {
	type st struct {
		Object map[string]any `json:"object,omitempty"`
		Array  []any          `json:"array,omitempty"`
		String string         `json:"string,omitempty"`
		Bool   bool           `json:"bool,omitempty"`
		Num    float32        `json:"num,omitempty"`
		Null   any            `json:"null,omitempty"`
	}

	cv("by default", func() {
		s := st{}
		v, err := Import(s)
		so(err, isNil)

		t.Logf("default omitempty: %s", v.MustMarshalString(OptSetSequence()))
		so(v.Len(), eq, 0)
	})

	cv("ignore omitempty", func() {
		s := st{}
		v, err := Import(s, OptIgnoreOmitempty())
		so(err, isNil)

		t.Logf("after ignoring omitempty: %s", v.MustMarshalString(OptSetSequence()))
		so(v.Len(), eq, 6)
	})

	cv("ignore omitempty in array", func() {
		arr := []st{{}, {}}
		v, err := Import(arr, OptIgnoreOmitempty())
		so(err, isNil)

		t.Logf("after ignoring omitempty: %s", v.MustMarshalString(OptSetSequence()))
		so(v.Len(), eq, 2)
		so(v.MustGet(0).Len(), eq, 6)
		so(v.MustGet(1).Len(), eq, 6)
	})

	cv("ignore omitempty in map", func() {
		arr := map[string]st{
			"a": {},
			"b": {},
		}
		v, err := Import(arr, OptIgnoreOmitempty())
		so(err, isNil)

		t.Logf("after ignoring omitempty: %s", v.MustMarshalString(OptSetSequence()))
		so(v.Len(), eq, 2)
		so(v.MustGet("a").Len(), eq, 6)
		so(v.MustGet("b").Len(), eq, 6)
	})
}

// https://github.com/Andrew-M-C/go.jsonvalue/issues/29
func testOption_Issue29(*testing.T) {
	cv("OptSetSequence in V in struct", func() {
		st := issue29Struct{}
		st.Ext = NewObject()

		st.Ext.MustSetString("1111").At("A")
		st.Ext.MustSetString("2222").At("B")

		const expected = `{"ext":{"A":"1111","B":"2222"}}`
		for i := 0; i < 2000; i++ {
			b, _ := json.Marshal(st)
			so(string(b), eq, expected)
		}
	})
}

type issue29Struct struct {
	Ext *V `json:"ext"`
}

func (s issue29Struct) MarshalJSON() ([]byte, error) {
	w := issue29StructWrapper(s)
	return New(w).Marshal(OptSetSequence())
}

type issue29StructWrapper issue29Struct
