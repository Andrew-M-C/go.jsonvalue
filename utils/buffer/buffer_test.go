package buffer_test

import (
	"testing"

	"github.com/Andrew-M-C/go.jsonvalue/utils/buffer"
	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So
	eq = convey.ShouldEqual
	ne = convey.ShouldNotEqual
)

func TestBuffer(t *testing.T) {
	cv("general", t, func() { testGeneral(t) })
	cv("JSON-like text", t, func() { testJSONLikeText(t) })
}

func testGeneral(t *testing.T) {
	buf := buffer.NewBuffer()
	so(buf, ne, nil)

	buf.WriteString("12345678")
	buf.WriteByte('A')
	buf.WriteRune('一')

	b := buf.Bytes()
	so(string(b), eq, "12345678A一")
}

func testJSONLikeText(t *testing.T) {
	buf := buffer.NewBuffer()
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
