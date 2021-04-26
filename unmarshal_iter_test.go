package jsonvalue

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func test(t *testing.T, scene string, f func(*testing.T)) {
	if t.Failed() {
		return
	}
	Convey(scene, t, func() {
		f(t)
	})
}

func TestIter(t *testing.T) {
	test(t, "iter.memcpy", testIter_memcpy)
	test(t, "iter.assignWideRune", testIter_assignWideRune)
	test(t, "iter.parseStrFromBytesBackward", testIter_generalStringUnmarshal)
	test(t, "iter.parseStrFromBytesForward and Backward", testIter_parseStrFromBytesBackwardForward)
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
		len += it.assignWideRune(len, r)
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

func printBytes(t *testing.T, b []byte) {
	if len(b) == 0 {
		t.Log("nil bytes")
		return
	}
	repeat := (len(b)-1)/10 + 1
	s := strings.Repeat(" 0 1 2 3 4 5 6 7 8 9", repeat)
	t.Log("")
	t.Log(s[:len(b)*2])
	t.Log(hex.EncodeToString(b))
}

func testIter_generalStringUnmarshal(t *testing.T) {
	b := []byte("hello, 世界")
	it := iter{b: b}

	printBytes(t, b)

	le, err := it.parseStrFromBytesBackward(0, len(b))
	t.Logf("le = %d, err = %v", le, nil)
	So(err, ShouldBeNil)
	So(le, ShouldEqual, len(b))
}

func testIter_parseStrFromBytesBackwardForward(t *testing.T) {
	// orig := "你👨‍👩‍👧‍👧你"
	orig := fmt.Sprintf(
		"%c%c%c%c%c%c%c%c%c",
		0x2F804, 0x1F468, 0x200D, 0x1F469, 0x200D, 0x1F467, 0x200D, 0x1F467, 0x4F60,
	)
	v := NewString(orig)

	Convey("backward", func() {
		raw := v.MustMarshal()
		t.Log("raw data:")
		printBytes(t, raw)

		it := iter{
			b: raw[1 : len(raw)-1],
		}

		t.Log("raw string:")
		printBytes(t, it.b)
		le, err := it.parseStrFromBytesBackward(0, len(it.b))
		So(err, ShouldBeNil)
		So(le, ShouldBeGreaterThan, 0)

		s := string(it.b[:le])
		t.Logf("Got len: %d", le)
		t.Logf("got string: %s", s)

		buff := bytes.Buffer{}
		for _, r := range s {
			buff.WriteString(fmt.Sprintf("0x%04x ", r))
		}
		t.Logf(buff.String())

		So(s, ShouldEqual, orig)
	})

	Convey("forward", func() {
		raw := v.MustMarshal()
		t.Log("raw data:")
		printBytes(t, raw)

		it := iter{b: raw}

		le, end, err := it.parseStrFromBytesForwardWithQuote(0)
		So(err, ShouldBeNil)
		So(le, ShouldBeGreaterThan, 0)

		t.Logf("Got len: %d, end: %d", le, end)

		s := string(it.b[1 : 1+le])
		t.Logf("got string: %s", s)

		buff := bytes.Buffer{}
		for _, r := range s {
			buff.WriteString(fmt.Sprintf("0x%04x ", r))
		}
		t.Logf(buff.String())

		So(s, ShouldEqual, orig)
	})
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

	end := len(raw)
	end, err := it.searchObjEnd(offset, end)
	t.Logf("end %d, err %v", end, err)
	So(err, ShouldBeNil)
	So(raw[end-1], ShouldEqual, '}')

	offset, reachEnd = it.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '[')

	end, err = it.searchArrEnd(offset, end-1)
	t.Logf("end %d, err %v", end, err)
	So(err, ShouldBeNil)
	So(raw[end-1], ShouldEqual, ']')

	offset, reachEnd = it.skipBlanks(offset + 1)
	t.Logf("offset %d, reachEnd %v", offset, reachEnd)
	So(offset, ShouldNotBeZeroValue)
	So(reachEnd, ShouldBeFalse)
	So(raw[offset], ShouldEqual, '{')

	end, err = it.searchObjEnd(offset, end-1)
	t.Logf("end %d, err %v", end, err)
	So(err, ShouldBeNil)
	So(raw[end-1], ShouldEqual, '}')

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

		i64, u64, f64, _, _, end, reachEnd, err := it.parseNumber(0)
		t.Logf("i64 = %v, u64 = %v, f64 = %v", i64, u64, f64)
		t.Logf("end = %d, readnEnd = %v", end, reachEnd)
		t.Logf(string(b[:end]))
		So(err, ShouldBeNil)
		So(f64, ShouldEqual, -12345.6789)
		So(reachEnd, ShouldBeTrue)
	})

	Convey("reachEnd == false", func() {
		it := &iter{b: b}

		i64, u64, f64, _, _, end, reachEnd, err := it.parseNumber(0)
		t.Logf("i64 = %v, u64 = %v, f64 = %v", i64, u64, f64)
		t.Logf("end = %d, readnEnd = %v", end, reachEnd)
		t.Logf(string(b[:end]))
		So(err, ShouldBeNil)
		So(f64, ShouldEqual, -12345.6789)
		So(reachEnd, ShouldBeFalse)
	})
}
