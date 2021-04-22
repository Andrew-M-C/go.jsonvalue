package jsonvalue

import (
	"bytes"
	"encoding/hex"
	"fmt"
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

func TestUtf8Iter(t *testing.T) {
	test(t, "utf8Iter.memcpy", testUtf8Iter_memcpy)
	test(t, "utf8Iter.assignWideRune", testUtf8Iter_assignWideRune)
	test(t, "utf8Iter.parseStrFromBytes", testUtf8Iter_parseStrFromBytes)
}

func testUtf8Iter_memcpy(t *testing.T) {
	b := []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA}

	it := utf8Iter{
		b: b,
	}

	origByte := b[4]

	t.Logf("before: %s", hex.EncodeToString(b))
	it.memcpy(0, 4, 6)
	t.Logf("result: %s", hex.EncodeToString(b))

	So(b[0], ShouldEqual, origByte)
}

func testUtf8Iter_assignWideRune(t *testing.T) {
	b := make([]byte, 32)

	it := utf8Iter{
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

func testUtf8Iter_parseStrFromBytes(t *testing.T) {
	// orig := "你👨‍👩‍👧‍👧你"
	orig := fmt.Sprintf(
		"%c%c%c%c%c%c%c%c%c",
		0x2F804, 0x1F468, 0x200D, 0x1F469, 0x200D, 0x1F467, 0x200D, 0x1F467, 0x4F60,
	)
	v := NewString(orig)

	raw := v.MustMarshal()
	t.Logf("raw data: %s", raw)

	it := utf8Iter{
		b: raw[1 : len(raw)-1],
	}
	t.Logf("raw string: %s", it.b)
	le, err := it.parseStrFromBytes(0, len(it.b))
	So(err, ShouldBeNil)
	So(le, ShouldNotBeZeroValue)

	s := string(it.b[:le])
	t.Logf("got string: %s", s)

	buff := bytes.Buffer{}
	for _, r := range s {
		buff.WriteString(fmt.Sprintf("0x%04x ", r))
	}
	t.Logf(buff.String())

	So(s, ShouldEqual, orig)
}
