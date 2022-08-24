package jsonvalue

import (
	"encoding/hex"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func testIter(t *testing.T) {
	cv("u.memcpy", func() { testIterMemcpy(t) })
	cv("u.assignWideRune", func() { testIterAssignWideRune(t) })
	cv("u.character searching", func() { testIterChrSearching(t) })
	cv("u.testIter_parseNumber", func() { testIterParseNumber(t) })
}

func testIterMemcpy(t *testing.T) {
	b := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA}

	u := newUnmarshaler(b)

	origByte := b[4]

	t.Logf("before: %s", hex.EncodeToString(b))
	u.memcpy(0, 4, 6)
	t.Logf("result: %s", hex.EncodeToString(b))

	So(b[0], ShouldEqual, origByte)
}

func testIterAssignWideRune(t *testing.T) {
	b := make([]byte, 32)

	u := newUnmarshaler(b)

	len := 0

	append := func(r rune) {
		t.Logf("rune hex: %04x", r)
		len += u.assignASCIICodedRune(len, r)
		t.Logf("bytes: %v", hex.EncodeToString(b))
	}

	append('您')
	append('好')
	append('世')
	append('界')

	u.b[len] = '!'
	len++

	b = b[:len]
	So(string(b), ShouldEqual, "您好世界!")
}

func testIterChrSearching(t *testing.T) {
	raw := []byte("   {  [ {  } ]  }  ")
	t.Logf("")
	t.Logf(string(raw))
	t.Logf("01234567890123456789")

	u := newUnmarshaler(raw)

	offset, reachEnd := u.skipBlanks(0)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '{')

	offset, reachEnd = u.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '[')

	offset, reachEnd = u.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '{')

	offset, reachEnd = u.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '}')
}

func testIterParseNumber(t *testing.T) {
	b := []byte("-12345.6789  ")

	Convey("reachEnd == true", func() {
		u := newUnmarshaler(b[:11])

		v, end, reachEnd, err := u.parseNumber(0)
		t.Logf("i64 = %v, u64 = %v, f64 = %v", v.num.i64, v.num.u64, v.num.f64)
		t.Logf("end = %d, readnEnd = %v", end, reachEnd)
		t.Logf(string(b[:end]))
		So(err, ShouldBeNil)
		So(v.num.f64, ShouldEqual, -12345.6789)
		So(reachEnd, ShouldBeTrue)
	})

	Convey("reachEnd == false", func() {
		u := newUnmarshaler(b)

		v, end, reachEnd, err := u.parseNumber(0)
		So(err, ShouldBeNil)
		So(v.num.f64, ShouldEqual, -12345.6789)
		So(reachEnd, ShouldBeFalse)
		t.Logf("i64 = %v, u64 = %v, f64 = %v", v.num.i64, v.num.u64, v.num.f64)
		t.Logf("end = %d, readnEnd = %v", end, reachEnd)
		t.Logf(string(b[:end]))
	})
}
