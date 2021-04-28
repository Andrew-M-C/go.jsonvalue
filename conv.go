package jsonvalue

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func parseUint(b []byte) (uint64, error) {
	return strconv.ParseUint(unsafeBtoS(b), 10, 64)
}

func parseInt(b []byte) (int64, error) {
	return strconv.ParseInt(unsafeBtoS(b), 10, 64)
}

func parseFloat(b []byte) (float64, error) {
	return strconv.ParseFloat(unsafeBtoS(b), 64)
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
	if r <= 0x7F {
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
