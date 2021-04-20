package jsonvalue

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

func parseUint(b []byte) (uint64, error) {
	return strconv.ParseUint(string(b), 10, 64)
}

func parseInt(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}

func parseFloat(b []byte) (float64, error) {
	return strconv.ParseFloat(string(b), 64)
}

// reference: https://golang.org/src/encoding/json/decode.go, func unquote()
func parseString(b []byte) (string, []byte, error) {
	firstQuote := bytes.Index(b, []byte{'"'})
	lastQuote := bytes.LastIndex(b, []byte{'"'})
	if firstQuote < 0 || firstQuote == lastQuote {
		return "", nil, ErrIllegalString
	}

	b = b[firstQuote : lastQuote-firstQuote+1]

	t, ok := unquoteBytes(b)
	if !ok {
		return "", nil, fmt.Errorf("invalid string '%s'", string(b))
	}
	return string(t), b, nil
}

func unquoteBytes(s []byte) (t []byte, ok bool) {
	if len(s) < 2 || s[0] != '"' || s[len(s)-1] != '"' {
		return
	}
	s = s[1 : len(s)-1]

	// Check for unusual characters. If there are none,
	// then no unquoting is needed, so return a slice of the
	// original bytes.
	r := 0
	for r < len(s) {
		c := s[r]
		if c == '\\' || c == '"' || c < ' ' {
			break
		}
		if c < utf8.RuneSelf {
			r++
			continue
		}
		rr, size := utf8.DecodeRune(s[r:])
		if rr == utf8.RuneError && size == 1 {
			break
		}
		r += size
	}
	if r == len(s) {
		return s, true
	}

	b := make([]byte, len(s)+2*utf8.UTFMax)
	w := copy(b, s[0:r])
	for r < len(s) {
		// Out of room? Can only happen if s is full of
		// malformed UTF-8 and we're replacing each
		// byte with RuneError.
		if w >= len(b)-2*utf8.UTFMax {
			nb := make([]byte, (len(b)+utf8.UTFMax)*2)
			copy(nb, b[0:w])
			b = nb
		}
		switch c := s[r]; {
		case c == '\\':
			r++
			if r >= len(s) {
				return
			}
			switch s[r] {
			default:
				return
			case '"', '\\', '/', '\'':
				b[w] = s[r]
				r++
				w++
			case 'b':
				b[w] = '\b'
				r++
				w++
			case 'f':
				b[w] = '\f'
				r++
				w++
			case 'n':
				b[w] = '\n'
				r++
				w++
			case 'r':
				b[w] = '\r'
				r++
				w++
			case 't':
				b[w] = '\t'
				r++
				w++
			case 'u':
				r--
				rr := getu4(s[r:])
				if rr < 0 {
					return
				}
				r += 6
				if utf16.IsSurrogate(rr) {
					rr1 := getu4(s[r:])
					if dec := utf16.DecodeRune(rr, rr1); dec != unicode.ReplacementChar {
						// A valid pair; consume.
						r += 6
						w += utf8.EncodeRune(b[w:], dec)
						break
					}
					// Invalid surrogate; fall back to replacement rune.
					rr = unicode.ReplacementChar
				}
				w += utf8.EncodeRune(b[w:], rr)
			}

		// Quote, control characters are invalid.
		case c == '"', c < ' ':
			return

		// ASCII
		case c < utf8.RuneSelf:
			b[w] = c
			r++
			w++

		// Coerce to well-formed UTF-8.
		default:
			rr, size := utf8.DecodeRune(s[r:])
			r += size
			w += utf8.EncodeRune(b[w:], rr)
		}
	}
	return b[0:w], true
}

func getu4(s []byte) rune {
	if len(s) < 6 || s[0] != '\\' || s[1] != 'u' {
		return -1
	}
	var r rune
	for _, c := range s[2:6] {
		switch {
		case '0' <= c && c <= '9':
			c = c - '0'
		case 'a' <= c && c <= 'f':
			c = c - 'a' + 10
		case 'A' <= c && c <= 'F':
			c = c - 'A' + 10
		default:
			return -1
		}
		r = r*16 + rune(c)
	}
	return r
}

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
func escapeUnicodeToBuff(buf *bytes.Buffer, r rune) {
	if r <= '\u0127' {
		buf.WriteRune(r)
		return
	}
	if r <= '\uffff' {
		buf.WriteString(fmt.Sprintf("\\u%04X", r))
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
	buf.WriteString(fmt.Sprintf("\\u%04X", hi+0xD800))
	buf.WriteString(fmt.Sprintf("\\u%04X", lo+0xDC00))
}

var (
	escDoubleQuote = []byte{'\\', '"'}
	escSlash       = []byte{'\\', '/'}
	escBaskslash   = []byte{'\\', '\\'}
	escBaskspace   = []byte{'\\', 'b'}
	escVertTab     = []byte{'\\', 'f'}
	escTab         = []byte{'\\', 't'}
	escNewLine     = []byte{'\\', 'n'}
	escReturn      = []byte{'\\', 'r'}
	escLeftAngle   = []byte{'\\', 'u', '0', '0', '3', 'C'}
	escRightAngle  = []byte{'\\', 'u', '0', '0', '3', 'E'}
	escAnd         = []byte{'\\', 'u', '0', '0', '2', '6'}
	escPercent     = []byte{'\\', 'u', '0', '0', '2', '5'}
)

func escapeStringToBuff(s string, buf *bytes.Buffer) {
	for _, chr := range s {
		switch chr {
		case '"':
			// buf.WriteString("\\\"")
			buf.Write(escDoubleQuote)
		case '/':
			// buf.WriteString("\\/")
			buf.Write(escSlash)
		case '\\':
			// buf.WriteString("\\\\")
			buf.Write(escBaskslash)
		case '\b':
			// buf.WriteString("\\b")
			buf.Write(escBaskspace)
		case '\f':
			// buf.WriteString("\\f")
			buf.Write(escVertTab)
		case '\t':
			// buf.WriteString("\\t")
			buf.Write(escTab)
		case '\n':
			// buf.WriteString("\\n")
			buf.Write(escNewLine)
		case '\r':
			// buf.WriteString("\\r")
			buf.Write(escReturn)
		case '<':
			// buf.WriteString("\\u003C")
			buf.Write(escLeftAngle)
		case '>':
			// buf.WriteString("\\u003E")
			buf.Write(escRightAngle)
		case '&':
			// buf.WriteString("\\u0026")
			buf.Write(escAnd)
		case '%': // not standard JSON encoding
			// buf.WriteString("\\u0025")
			buf.Write(escPercent)
		default:
			escapeUnicodeToBuff(buf, chr)
		}
	}
}

func intfToInt(v interface{}) (u int, err error) {
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

// func intfToInt64(v interface{}) (i int64, err error) {
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

func intfToString(v interface{}) (s string, err error) {
	switch str := v.(type) {
	case string:
		return str, nil
	default:
		return "", fmt.Errorf("%s is not a string", reflect.TypeOf(v).String())
	}
}

// func intfToJsonvalue(v interface{}) (j *V, err error) {
// 	switch v.(type) {
// 	case *V:
// 		j = v.(*V)
// 	default:
// 		err = fmt.Errorf("%s is not a *jsonvalue.V type", reflect.TypeOf(v).String())
// 	}

// 	return
// }
