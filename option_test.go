package jsonvalue

import (
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
}

func testOptionOverwriting(t *testing.T) {
	v := NewObject(M{
		"slash": "/",
	})

	expect := `{"slash":"/"}`
	s := v.MustMarshalString()
	so(s, eq, expect)
}

func testOptionReset(t *testing.T) {
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

func testOption_OptSetSequence(t *testing.T) {
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
