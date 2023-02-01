package jsonvalue

import (
	"testing"
)

func testBuffer(t *testing.T) {
	cv("general", func() { testGeneral(t) })
	cv("JSON-like text", func() { testJSONLikeText(t) })
}

func testGeneral(t *testing.T) {
	buf := NewBuffer()
	so(buf, ne, nil)

	buf.WriteString("12345678")
	buf.WriteByte('A')
	buf.WriteRune('一')

	b := buf.Bytes()
	so(string(b), eq, "12345678A一")
}

func testJSONLikeText(t *testing.T) {
	buf := NewBuffer()
	buf.WriteByte('{')

	buf.WriteByte('"')
	buf.WriteString("message")
	buf.WriteByte('"')

	buf.WriteByte(':')
	buf.WriteByte('"')
	buf.WriteString("Hello, world!")
	buf.WriteByte('"')
	buf.WriteByte('}')

	b := buf.Bytes()
	so(string(b), eq, `{"message":"Hello, world!"}`)
}
