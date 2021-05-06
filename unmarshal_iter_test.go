package jsonvalue

import (
	"encoding/hex"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIter(t *testing.T) {
	test(t, "iter.memcpy", testIter_memcpy)
	test(t, "iter.assignWideRune", testIter_assignWideRune)
	test(t, "iter.character searching", testIter_chrSearching)
	test(t, "iter.testIter_parseNumber", testIter_parseNumber)
}

func testIter_memcpy(t *testing.T) {
	b := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA}

	it := iter{b: b}

	origByte := b[4]

	t.Logf("before: %s", hex.EncodeToString(b))
	it.memcpy(0, 4, 6)
	t.Logf("result: %s", hex.EncodeToString(b))

	So(b[0], ShouldEqual, origByte)
}

func testIter_assignWideRune(t *testing.T) {
	b := make([]byte, 32)

	it := iter{
		b: b,
	}

	len := 0

	append := func(r rune) {
		t.Logf("rune hex: %04x", r)
		len += it.assignAsciiCodedRune(len, r)
		t.Logf("bytes: %v", hex.EncodeToString(b))
	}

	append('您')
	append('好')
	append('世')
	append('界')

	it.b[len] = '!'
	len++

	b = b[:len]
	So(string(b), ShouldEqual, "您好世界!")
}

func testIter_chrSearching(t *testing.T) {
	raw := []byte("   {  [ {  } ]  }  ")
	t.Logf("")
	t.Logf(string(raw))
	t.Logf("01234567890123456789")

	it := iter{b: raw}

	offset, reachEnd := it.skipBlanks(0)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '{')

	offset, reachEnd = it.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '[')

	offset, reachEnd = it.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '{')

	offset, reachEnd = it.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '}')
}

func testIter_parseNumber(t *testing.T) {
	b := []byte("-12345.6789  ")

	Convey("reachEnd == true", func() {
		it := &iter{b: b[:11]}

		v, end, reachEnd, err := it.parseNumber(0)
		t.Logf("i64 = %v, u64 = %v, f64 = %v", v.num.i64, v.num.u64, v.num.f64)
		t.Logf("end = %d, readnEnd = %v", end, reachEnd)
		t.Logf(string(b[:end]))
		So(err, ShouldBeNil)
		So(v.num.f64, ShouldEqual, -12345.6789)
		So(reachEnd, ShouldBeTrue)
	})

	Convey("reachEnd == false", func() {
		it := &iter{b: b}

		v, end, reachEnd, err := it.parseNumber(0)
		t.Logf("i64 = %v, u64 = %v, f64 = %v", v.num.i64, v.num.u64, v.num.f64)
		t.Logf("end = %d, readnEnd = %v", end, reachEnd)
		t.Logf(string(b[:end]))
		So(err, ShouldBeNil)
		So(v.num.f64, ShouldEqual, -12345.6789)
		So(reachEnd, ShouldBeFalse)
	})
}
