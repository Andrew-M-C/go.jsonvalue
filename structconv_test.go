package jsonvalue

import (
	"testing"
)

func testStructConv(t *testing.T) {
	cv("export to string", func() { testExportString(t) })
	cv("export to int", func() { testExportInt(t) })
	cv("export to float", func() { testExportFloat(t) })
	cv("export to bool", func() { testExportBool(t) })
	cv("misc import", func() { testImport(t) })
}

func testExportString(t *testing.T) {
	const S = "Hello, jsonvalue!"
	v := NewString(S)

	str := ""
	err := v.Export(str)
	so(err, isErr)

	err = v.Export(&str)
	so(err, isNil)

	so(str, eq, S)

	bol := true
	err = v.Export(&bol)
	so(err, isErr)

	v = &V{}
	err = v.Export(nil)
	so(err, isErr)
}

func testExportInt(t *testing.T) {
	const positive = 123454321
	const negative = -987656789

	n1 := NewInt(positive)

	var i int
	var u uint
	var i32 int32
	var u32 uint32

	err := n1.Export(&i)
	so(err, isNil)
	so(i, eq, positive)

	err = n1.Export(&u)
	so(err, isNil)
	so(u, eq, positive)

	err = n1.Export(&i32)
	so(err, isNil)
	so(i32, eq, positive)

	err = n1.Export(&u32)
	so(err, isNil)
	so(u32, eq, positive)

	// --------

	n2 := NewInt(negative)

	err = n2.Export(&i)
	so(err, isNil)
	so(i, eq, negative)

	err = n2.Export(&i32)
	so(err, isNil)
	so(i32, eq, negative)

	// --------

	bol := true
	err = n1.Export(&bol)
	so(err, isErr)
}

func testExportFloat(t *testing.T) {
	const F = 12345.4321

	n := NewFloat64(F)

	var f32 float32
	var f64 float64

	err := n.Export(&f32)
	so(err, isNil)
	so(f32, eq, F)

	err = n.Export(&f64)
	so(err, isNil)
	so(f64, eq, F)

	// --------

	bol := true
	err = n.Export(&bol)
	so(err, isErr)
}

func testExportBool(t *testing.T) {
	v := NewBool(true)
	b := false

	err := v.Export(b)
	so(err, isErr)

	err = v.Export(&b)
	so(err, isNil)

	so(b, isTrue)

	str := ""
	err = v.Export(&str)
	so(err, isErr)
}

func testImport(t *testing.T) {
	cv("integers", func() {

		params := []interface{}{
			int(1),
			uint(2),
			int8(3),
			uint8(4),
			int16(5),
			uint16(6),
			int32(7),
			uint32(8),
			int64(9),
			uint64(10),
		}

		for i, p := range params {
			v, err := Import(p)
			so(err, isNil)
			so(v.ValueType(), eq, Number)
			so(v.Int(), eq, i+1)
		}
	})

	cv("string", func() {
		s := "hello"
		v, err := Import(s)
		so(err, isNil)
		so(v.ValueType(), eq, String)
		so(v.String(), eq, s)
	})

	cv("object", func() {
		type thing struct {
			String string `json:"str"`
		}
		th := thing{
			String: "world",
		}

		v, err := Import(&th)
		so(err, isNil)
		so(v.ValueType(), eq, Object)

		s, err := v.GetString("str")
		so(err, isNil)
		so(s, eq, th.String)
	})

	cv("error", func() {
		f := func() bool {
			return false
		}
		v, err := Import(f)
		so(err, isErr)
		so(v, notNil)
		so(v.ValueType(), eq, NotExist)
	})
}
