package jsonvalue

import (
	"fmt"
	"reflect"

	"github.com/Andrew-M-C/go.jsonvalue/internal/buffer"
)

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// reference:
// - [UTF-16](https://zh.wikipedia.org/zh-cn/UTF-16)
// - [JavaScript has a Unicode problem](https://mathiasbynens.be/notes/javascript-unicode)
// - [Meaning of escaped unicode characters in JSON](https://stackoverflow.com/questions/21995410/meaning-of-escaped-unicode-characters-in-json)
func escapeGreaterUnicodeToBuffByUTF16(r rune, buf buffer.Buffer) {
	if r <= '\uffff' {
		_, _ = buf.WriteString(fmt.Sprintf("\\u%04X", r))
		return
	}
	// if r > 0x10FFFF {
	// 	// invalid unicode
	// 	buf.WriteRune(r)
	// 	return
	// }

	r = r - 0x10000
	lo := r & 0x003FF
	hi := (r & 0xFFC00) >> 10
	_, _ = buf.WriteString(fmt.Sprintf("\\u%04X", hi+0xD800))
	_, _ = buf.WriteString(fmt.Sprintf("\\u%04X", lo+0xDC00))
}

func escapeGreaterUnicodeToBuffByUTF8(r rune, buf buffer.Buffer) {
	// Comments below are copied from encoding/json:
	//
	// U+2028 is LINE SEPARATOR.
	// U+2029 is PARAGRAPH SEPARATOR.
	// They are both technically valid characters in JSON strings,
	// but don't work in JSONP, which has to be evaluated as JavaScript,
	// and can lead to security holes there. It is valid JSON to
	// escape them, so we do so unconditionally.
	// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
	if r == '\u2028' || r == '\u2029' {
		escapeGreaterUnicodeToBuffByUTF16(r, buf)
	} else {
		_, _ = buf.WriteRune(r)
	}
}

func escapeNothing(b byte, buf buffer.Buffer) {
	_ = buf.WriteByte(b)
}

func escAsciiControlChar(b byte, buf buffer.Buffer) {
	upper := b >> 4
	lower := b & 0x0F

	writeChar := func(c byte) {
		if c < 0xA {
			_ = buf.WriteByte('0' + c)
		} else {
			_ = buf.WriteByte('A' + (c - 0xA))
		}
	}

	_, _ = buf.Write([]byte{'\\', 'u', '0', '0'})
	writeChar(upper)
	writeChar(lower)
}

func escDoubleQuote(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', '"'})
}

func escSlash(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', '/'})
}

func escBackslash(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', '\\'})
}

func escBackspace(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 'b'})
}

func escVertTab(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 'f'})
}

func escTab(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 't'})
}

func escNewLine(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 'n'})
}

func escReturn(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 'r'})
}

func escLeftAngle(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 'u', '0', '0', '3', 'C'})
}

func escRightAngle(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 'u', '0', '0', '3', 'E'})
}

func escAnd(_ byte, buf buffer.Buffer) {
	_, _ = buf.Write([]byte{'\\', 'u', '0', '0', '2', '6'})
}

// func escPercent(_ byte, buf buffer.Buffer) {
// 	buf.Write([]byte{'\\', 'u', '0', '0', '2', '5'})
// }

func escapeStringToBuff(s string, buf buffer.Buffer, opt *Opt) {
	for _, r := range s {
		if r <= 0x7F {
			b := byte(r)
			opt.asciiCharEscapingFunc[b](b, buf)
		} else {
			opt.unicodeEscapingFunc(r, buf)
		}
	}
}

func intfToInt(v any) (u int, err error) {
	switch v := v.(type) {
	case int:
		u = v
	case uint:
		u = int(v)
	case int64:
		u = int(v)
	case uint64:
		u = int(v)
	case int32:
		u = int(v)
	case uint32:
		u = int(v)
	case int16:
		u = int(v)
	case uint16:
		u = int(v)
	case int8:
		u = int(v)
	case uint8:
		u = int(v)
	default:
		err = fmt.Errorf("%s is not a number", reflect.TypeOf(v).String())
	}

	return
}

// func intfToInt64(v any) (i int64, err error) {
// 	switch v.(type) {
// 	case int:
// 		i = int64(v.(int))
// 	case uint:
// 		i = int64(v.(uint))
// 	case int64:
// 		i = int64(v.(int64))
// 	case uint64:
// 		i = int64(v.(uint64))
// 	case int32:
// 		i = int64(v.(int32))
// 	case uint32:
// 		i = int64(v.(uint32))
// 	case int16:
// 		i = int64(v.(int16))
// 	case uint16:
// 		i = int64(v.(uint16))
// 	case int8:
// 		i = int64(v.(int8))
// 	case uint8:
// 		i = int64(v.(uint8))
// 	default:
// 		err = fmt.Errorf("%s is not a number", reflect.TypeOf(v).String())
// 	}

// 	return
// }

func intfToString(v any) (s string, err error) {
	switch str := v.(type) {
	case string:
		return str, nil
	default:
		return "", fmt.Errorf("%s is not a string", reflect.TypeOf(v).String())
	}
}

// func intfToJsonvalue(v any) (j *V, err error) {
// 	switch v.(type) {
// 	case *V:
// 		j = v.(*V)
// 	default:
// 		err = fmt.Errorf("%s is not a *jsonvalue.V type", reflect.TypeOf(v).String())
// 	}

// 	return
// }
