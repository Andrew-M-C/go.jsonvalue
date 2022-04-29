package jsonvalue

import (
	"testing"
)

// This function used to ovwrwrite system default marshal options
func init() {
	SetDefaultMarshalOptions(OptEscapeSlash(false))
}

func testOption(t *testing.T) {
	cv("test default option overwriting", func() { testOptionOverwriting(t) })
	cv("test reset marshal options", func() { testOptionReset(t) })
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
