package jsonvalue

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Andrew-M-C/go.jsonvalue/internal/unsafe"
)

// ================ OUTER UNMARSHAL ================

// unmarshalWithIter parse bytes with unknown value type.
func unmarshalWithIter(p pool, it iter, offset int) (v *V, err error) {
	end := len(it)
	offset, reachEnd := it.skipBlanks(offset)
	if reachEnd {
		return &V{}, fmt.Errorf("%w, cannot find any symbol characters found", ErrRawBytesUnrecignized)
	}

	chr := it[offset]
	switch chr {
	case '{':
		v, offset, err = unmarshalObjectWithIterUnknownEnd(p, it, offset, end)

	case '[':
		v, offset, err = unmarshalArrayWithIterUnknownEnd(p, it, offset, end)

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		var n *V
		n, offset, _, err = it.parseNumber(p, offset)
		if err == nil {
			v = n
		}

	case '"':
		var sectLenWithoutQuote int
		var sectEnd int
		sectLenWithoutQuote, sectEnd, err = it.parseStrFromBytesForwardWithQuote(offset)
		if err == nil {
			v, err = NewString(unsafe.BtoS(it[offset+1:offset+1+sectLenWithoutQuote])), nil
			offset = sectEnd
		}

	case 't':
		offset, err = it.parseTrue(offset)
		if err == nil {
			v = NewBool(true)
		}

	case 'f':
		offset, err = it.parseFalse(offset)
		if err == nil {
			v = NewBool(false)
		}

	case 'n':
		offset, err = it.parseNull(offset)
		if err == nil {
			v = NewNull()
		}

	default:
		return &V{}, fmt.Errorf("%w, invalid character \\u%04X at Position %d", ErrRawBytesUnrecignized, chr, offset)
	}

	if err != nil {
		return &V{}, err
	}

	if offset, reachEnd = it.skipBlanks(offset, end); !reachEnd {
		return &V{}, fmt.Errorf("%w, unnecessary trailing data remains at Position %d", ErrRawBytesUnrecignized, offset)
	}

	return v, nil
}

// unmarshalArrayWithIterUnknownEnd is similar with unmarshalArrayWithIter, though should start with '[',
// but it does not known where its ']' is
func unmarshalArrayWithIterUnknownEnd(p pool, it iter, offset, right int) (_ *V, end int, err error) {
	offset++
	arr := newArray(p)

	reachEnd := false

	for offset < right {
		// search for ending ']'
		offset, reachEnd = it.skipBlanks(offset, right)
		if reachEnd {
			// ']' not found
			return nil, -1, fmt.Errorf("%w, cannot find ']'", ErrNotArrayValue)
		}

		chr := it[offset]
		switch chr {
		case ']':
			return arr, offset + 1, nil

		case ',':
			offset++

		case '{':
			v, sectEnd, err := unmarshalObjectWithIterUnknownEnd(p, it, offset, right)
			if err != nil {
				return nil, -1, err
			}
			arr.appendToArr(v)
			offset = sectEnd

		case '[':
			v, sectEnd, err := unmarshalArrayWithIterUnknownEnd(p, it, offset, right)
			if err != nil {
				return nil, -1, err
			}
			arr.appendToArr(v)
			offset = sectEnd

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			var v *V
			v, sectEnd, _, err := it.parseNumber(p, offset)
			if err != nil {
				return nil, -1, err
			}
			arr.appendToArr(v)
			offset = sectEnd

		case '"':
			sectLenWithoutQuote, sectEnd, err := it.parseStrFromBytesForwardWithQuote(offset)
			if err != nil {
				return nil, -1, err
			}
			v := NewString(unsafe.BtoS(it[offset+1 : offset+1+sectLenWithoutQuote]))
			arr.appendToArr(v)
			offset = sectEnd

		case 't':
			sectEnd, err := it.parseTrue(offset)
			if err != nil {
				return nil, -1, err
			}
			arr.appendToArr(NewBool(true))
			offset = sectEnd

		case 'f':
			sectEnd, err := it.parseFalse(offset)
			if err != nil {
				return nil, -1, err
			}
			arr.appendToArr(NewBool(false))
			offset = sectEnd

		case 'n':
			sectEnd, err := it.parseNull(offset)
			if err != nil {
				return nil, -1, err
			}
			arr.appendToArr(NewNull())
			offset = sectEnd

		default:
			return nil, -1, fmt.Errorf("%w, invalid character \\u%04X at Position %d", ErrRawBytesUnrecignized, chr, offset)
		}
	}

	return nil, -1, fmt.Errorf("%w, cannot find ']'", ErrNotArrayValue)
}

func (v *V) appendToArr(child *V) {
	if v.children.arr == nil {
		v.children.arr = make([]*V, 0, initialArrayCapacity)
	}
	v.children.arr = append(v.children.arr, child)
}

// unmarshalObjectWithIterUnknownEnd unmarshal object from raw bytes. it[offset] must be '{'
func unmarshalObjectWithIterUnknownEnd(p pool, it iter, offset, right int) (_ *V, end int, err error) {
	offset++
	obj := newObject(p)

	keyStart, keyEnd := 0, 0
	colonFound := false

	reachEnd := false

	keyNotFoundErr := func() error {
		if keyEnd == 0 {
			return fmt.Errorf(
				"%w, missing key for another value at Position %d", ErrNotObjectValue, offset,
			)
		}
		if !colonFound {
			return fmt.Errorf(
				"%w, missing colon for key at Position %d", ErrNotObjectValue, offset,
			)
		}
		return nil
	}

	valNotFoundErr := func() error {
		if keyEnd > 0 {
			return fmt.Errorf(
				"%w, missing value for key '%s' at Position %d",
				ErrNotObjectValue, unsafe.BtoS(it[keyStart:keyEnd]), keyStart,
			)
		}
		return nil
	}

	for offset < right {
		offset, reachEnd = it.skipBlanks(offset, right)
		if reachEnd {
			// '}' not found
			return nil, -1, fmt.Errorf("%w, cannot find '}'", ErrNotObjectValue)
		}

		chr := it[offset]
		switch chr {
		case '}':
			if err = valNotFoundErr(); err != nil {
				return nil, -1, err
			}
			return obj, offset + 1, nil

		case ',':
			if err = valNotFoundErr(); err != nil {
				return nil, -1, err
			}
			offset++
			// continue

		case ':':
			if colonFound {
				return nil, -1, fmt.Errorf("%w, duplicate colon at Position %d", ErrNotObjectValue, keyStart)
			}
			colonFound = true
			if err = keyNotFoundErr(); err != nil {
				return nil, -1, err
			}
			offset++
			// continue

		case '{':
			if err = keyNotFoundErr(); err != nil {
				return nil, -1, err
			}
			v, sectEnd, err := unmarshalObjectWithIterUnknownEnd(p, it, offset, right)
			if err != nil {
				return nil, -1, err
			}
			obj.setToObjectChildren(unsafe.BtoS(it[keyStart:keyEnd]), v)
			keyEnd, colonFound = 0, false
			offset = sectEnd

		case '[':
			if err = keyNotFoundErr(); err != nil {
				return nil, -1, err
			}
			v, sectEnd, err := unmarshalArrayWithIterUnknownEnd(p, it, offset, right)
			if err != nil {
				return nil, -1, err
			}
			obj.setToObjectChildren(unsafe.BtoS(it[keyStart:keyEnd]), v)
			keyEnd, colonFound = 0, false
			offset = sectEnd

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			if err = keyNotFoundErr(); err != nil {
				return nil, -1, err
			}
			var v *V
			v, sectEnd, _, err := it.parseNumber(p, offset)
			if err != nil {
				return nil, -1, err
			}
			obj.setToObjectChildren(unsafe.BtoS(it[keyStart:keyEnd]), v)
			keyEnd, colonFound = 0, false
			offset = sectEnd

		case '"':
			if keyEnd > 0 {
				// string value
				if !colonFound {
					return nil, -1, fmt.Errorf("%w, missing value for key '%s' at Position %d",
						ErrNotObjectValue, unsafe.BtoS(it[keyStart:keyEnd]), keyStart,
					)
				}
				sectLenWithoutQuote, sectEnd, err := it.parseStrFromBytesForwardWithQuote(offset)
				if err != nil {
					return nil, -1, err
				}
				v := NewString(unsafe.BtoS(it[offset+1 : offset+1+sectLenWithoutQuote]))
				obj.setToObjectChildren(unsafe.BtoS(it[keyStart:keyEnd]), v)
				keyEnd, colonFound = 0, false
				offset = sectEnd

			} else {
				// string key
				sectLenWithoutQuote, sectEnd, err := it.parseStrFromBytesForwardWithQuote(offset)
				if err != nil {
					return nil, -1, err
				}
				keyStart, keyEnd = offset+1, offset+1+sectLenWithoutQuote
				offset = sectEnd
			}

		case 't':
			if err = keyNotFoundErr(); err != nil {
				return nil, -1, err
			}
			sectEnd, err := it.parseTrue(offset)
			if err != nil {
				return nil, -1, err
			}
			obj.setToObjectChildren(unsafe.BtoS(it[keyStart:keyEnd]), NewBool(true))
			keyEnd, colonFound = 0, false
			offset = sectEnd

		case 'f':
			if err = keyNotFoundErr(); err != nil {
				return nil, -1, err
			}
			sectEnd, err := it.parseFalse(offset)
			if err != nil {
				return nil, -1, err
			}
			obj.setToObjectChildren(unsafe.BtoS(it[keyStart:keyEnd]), NewBool(false))
			keyEnd, colonFound = 0, false
			offset = sectEnd

		case 'n':
			if err = keyNotFoundErr(); err != nil {
				return nil, -1, err
			}
			sectEnd, err := it.parseNull(offset)
			if err != nil {
				return nil, -1, err
			}
			obj.setToObjectChildren(unsafe.BtoS(it[keyStart:keyEnd]), NewNull())
			keyEnd, colonFound = 0, false
			offset = sectEnd

		default:
			return nil, -1, fmt.Errorf("%w, invalid character \\u%04X at Position %d", ErrRawBytesUnrecignized, chr, offset)
		}

	}

	return nil, -1, fmt.Errorf("%w, cannot find '}'", ErrNotObjectValue)
}

// parseNumber parse a number string. Reference:
//
// - [ECMA-404 The JSON Data Interchange Standard](https://www.json.org/json-en.html)
func (v *V) parseNumber(p pool) (err error) {
	it := iter(v.srcByte)

	parsed, end, reachEnd, err := it.parseNumber(p, 0)
	if err != nil {
		return err
	}
	if !reachEnd {
		return fmt.Errorf("invalid character: 0x%02x", v.srcByte[end])
	}

	*v = *parsed
	return nil
}

// ==== simple object parsing ====
func newFromNumber(p pool, b []byte) (ret *V, err error) {
	v := new(p, Number)
	v.srcByte = b
	return v, nil
}

// ================ GENERAL UNMARSHALING ================

// iter is used to iterate []byte text
type iter []byte

func (it iter) parseStrFromBytesForwardWithQuote(offset int) (sectLenWithoutQuote int, sectEnd int, err error) {
	offset++ // skip "
	end := len(it)
	sectEnd = offset

	shift := func(i *int, le int) {
		if end-*i < le {
			err = fmt.Errorf(
				"%w, expect at least %d remaining bytes, but got %d at Position %d",
				ErrIllegalString, end-*i, le, *i,
			)
			return
		}
		it.memcpy(sectEnd, *i, le)
		sectEnd += le
		*i += le
	}

	// iterate every byte
	for i := offset; i < end; {
		chr := it[i]

		// ACSII?
		if chr == '\\' {
			err = it.handleEscapeStart(&i, &sectEnd)
		} else if chr == '"' {
			// found end quote
			return sectEnd - offset, i + 1, nil
		} else if chr <= 0x7F {
			// shift(&i, 1)
			it[sectEnd] = it[i]
			i++
			sectEnd++
		} else if runeIdentifyingBytes2(chr) {
			shift(&i, 2)
		} else if runeIdentifyingBytes3(chr) {
			shift(&i, 3)
		} else if runeIdentifyingBytes4(chr) {
			shift(&i, 4)
		} else {
			err = fmt.Errorf("%w: illegal UTF8 string at Position %d", ErrIllegalString, i)
		}

		if err != nil {
			return -1, -1, err
		}
	}

	err = errors.New("ending double quote of a string is not found")
	return
}

func (it iter) handleEscapeStart(i *int, sectEnd *int) error {
	if len(it)-1-*i < 1 {
		return errors.New("escape symbol not followed by another character")
	}

	chr := it[*i+1]
	switch chr {
	default:
		return fmt.Errorf("unreconized character 0x%02X after escape symbol", chr)
	case '"', '\'', '/', '\\':
		it[*sectEnd] = chr
		*sectEnd++
		*i += 2
	case 'b':
		it[*sectEnd] = '\b'
		*sectEnd++
		*i += 2
	case 'f':
		it[*sectEnd] = '\f'
		*sectEnd++
		*i += 2
	case 'r':
		it[*sectEnd] = '\r'
		*sectEnd++
		*i += 2
	case 'n':
		it[*sectEnd] = '\n'
		*sectEnd++
		*i += 2
	case 't':
		it[*sectEnd] = '\t'
		*sectEnd++
		*i += 2
	case 'u':
		return it.handleEscapeUnicodeStartWithEnd(i, len(it)-1, sectEnd)
	}
	return nil
}

func (it iter) handleEscapeUnicodeStartWithEnd(i *int, end int, sectEnd *int) (err error) {
	if end-*i <= 5 {
		return errors.New("escape symbol not followed by another character")
	}

	b3 := chrToHex(it[*i+2], &err)
	b2 := chrToHex(it[*i+3], &err)
	b1 := chrToHex(it[*i+4], &err)
	b0 := chrToHex(it[*i+5], &err)
	if err != nil {
		return
	}

	r := (rune(b3) << 12) + (rune(b2) << 8) + (rune(b1) << 4) + rune(b0)

	// this rune is smaller than 0x10000
	if r <= 0xD7FF || r >= 0xE000 {
		le := it.assignASCIICodedRune(*sectEnd, r)
		*i += 6
		*sectEnd += le
		return nil
	}

	// reference: [JSON 序列化中的转义和 Unicode 编码](https://cloud.tencent.com/developer/article/1625557/)
	// should get another unicode-escaped character
	if end-*i <= 11 {
		return fmt.Errorf("insufficient UTF-16 data at offset %d", *i)
	}
	if it[*i+6] != '\\' || it[*i+7] != 'u' {
		return fmt.Errorf("expect unicode escape character at position %d but not", *i+6)
	}

	ex3 := chrToHex(it[*i+8], &err)
	ex2 := chrToHex(it[*i+9], &err)
	ex1 := chrToHex(it[*i+10], &err)
	ex0 := chrToHex(it[*i+11], &err)
	if err != nil {
		return
	}

	ex := (rune(ex3) << 12) + (rune(ex2) << 8) + (rune(ex1) << 4) + rune(ex0)
	if ex < 0xDC00 {
		return fmt.Errorf(
			"%w, expect second UTF-16 encoding but got 0x04%X at position %d",
			ErrIllegalString, r, *i+8,
		)
	}
	ex -= 0xDC00
	if ex > 0x03FF {
		return fmt.Errorf(
			"%w, expect second UTF-16 encoding but got 0x04%X at position %d",
			ErrIllegalString, r, *i+8,
		)
	}

	r = ((r - 0xD800) << 10) + ex + 0x10000

	le := it.assignASCIICodedRune(*sectEnd, r)
	*i += 12
	*sectEnd += le
	return nil
}

func chrToHex(chr byte, errOut *error) byte {
	if chr >= '0' && chr <= '9' {
		return chr - '0'
	}
	if chr >= 'A' && chr <= 'F' {
		return chr - 'A' + 10
	}
	if chr >= 'a' && chr <= 'f' {
		return chr - 'a' + 10
	}
	*errOut = fmt.Errorf("invalid unicode value character: %c", rune(chr))
	return 0
}

func (it iter) memcpy(dst, src, length int) {
	if dst == src {
		return
	}
	copy(it[dst:dst+length], it[src:src+length])
	// ptr := unsafe.Pointer(&it[0])
	// C.memcpy(
	// 	unsafe.Pointer(uintptr(ptr)+uintptr(dst)),
	// 	unsafe.Pointer(uintptr(ptr)+uintptr(src)),
	// 	C.size_t(length),
	// )
}

func (it iter) assignASCIICodedRune(dst int, r rune) (offset int) {
	// 0zzzzzzz ==>
	// 0zzzzzzz
	if r <= 0x7F {
		it[dst+0] = byte(r)
		return 1
	}

	// 00000yyy yyzzzzzz ==>
	// 110yyyyy 10zzzzzz
	if r <= 0x7FF {
		it[dst+0] = byte((r&0x7C0)>>6) + 0xC0
		it[dst+1] = byte((r&0x03F)>>0) + 0x80
		return 2
	}

	// xxxxyyyy yyzzzzzz ==>
	// 1110xxxx 10yyyyyy 10zzzzzz
	if r <= 0xFFFF {
		it[dst+0] = byte((r&0xF000)>>12) + 0xE0
		it[dst+1] = byte((r&0x0FC0)>>6) + 0x80
		it[dst+2] = byte((r&0x003F)>>0) + 0x80
		return 3
	}

	// 000wwwxx xxxxyyyy yyzzzzzz ==>
	// 11110www 10xxxxxx 10yyyyyy 10zzzzzz
	it[dst+0] = byte((r&0x1C0000)>>18) + 0xF0
	it[dst+1] = byte((r&0x03F000)>>12) + 0x80
	it[dst+2] = byte((r&0x000FC0)>>6) + 0x80
	it[dst+3] = byte((r&0x00003F)>>0) + 0x80
	return 4
}

func runeIdentifyingBytes2(chr byte) bool {
	return (chr & 0xE0) == 0xC0
}

func runeIdentifyingBytes3(chr byte) bool {
	return (chr & 0xF0) == 0xE0
}

func runeIdentifyingBytes4(chr byte) bool {
	return (chr & 0xF8) == 0xF0
}

func (it iter) parseTrue(offset int) (end int, err error) {
	if len(it)-offset < 4 {
		return -1, fmt.Errorf("%w, insufficient character from Position %d", ErrNotValidBoolValue, offset)
	}

	if it[offset] == 't' &&
		it[offset+1] == 'r' &&
		it[offset+2] == 'u' &&
		it[offset+3] == 'e' {
		return offset + 4, nil
	}

	return -1, fmt.Errorf("%w, not 'true' at Position %d", ErrNotValidBoolValue, offset)
}

func (it iter) parseFalse(offset int) (end int, err error) {
	if len(it)-offset < 5 {
		return -1, fmt.Errorf("%w, insufficient character from Position %d", ErrNotValidBoolValue, offset)
	}

	if it[offset] == 'f' &&
		it[offset+1] == 'a' &&
		it[offset+2] == 'l' &&
		it[offset+3] == 's' &&
		it[offset+4] == 'e' {
		return offset + 5, nil
	}

	return -1, fmt.Errorf("%w, not 'false' at Position %d", ErrNotValidBoolValue, offset)
}

func (it iter) parseNull(offset int) (end int, err error) {
	if len(it)-offset < 4 {
		return -1, fmt.Errorf("%w, insufficient character from Position %d", ErrNotValidNulllValue, offset)
	}

	if it[offset] == 'n' &&
		it[offset+1] == 'u' &&
		it[offset+2] == 'l' &&
		it[offset+3] == 'l' {
		return offset + 4, nil
	}

	return -1, fmt.Errorf("%w, not 'null' at Position %d", ErrNotValidBoolValue, offset)
}

// skipBlanks skip blank characters until end or reaching a non-blank characher
func (it iter) skipBlanks(offset int, endPos ...int) (newOffset int, reachEnd bool) {
	end := 0
	if len(endPos) > 0 {
		end = endPos[0]
	} else {
		end = len(it)
	}

	for offset < end {
		chr := it[offset]
		switch chr {
		case ' ', '\r', '\n', '\t', '\b':
			offset++ // continue
		default:
			return offset, false
		}
	}

	return end, true
}

// ================ FLOAT UNMARSHALING ================

// For state machine chart, please refer to ./img/parse_float_state_chart.drawio

func (it iter) parseNumber(
	p pool, offset int,
) (v *V, end int, reachEnd bool, err error) {

	idx := offset
	negative := false
	floated := false
	exponentGot := false
	dotGot := false
	intAfterDotGot := false
	integer := uint64(0)
	edgeFound := false

	// len(it)-idx means remain bytes

	for ; len(it)-idx > 0 && !edgeFound; idx++ {
		b := it[idx]

		switch b {
		default:
			edgeFound = true

		case '0':
			if idx == offset {
				// OK
			} else if exponentGot {
				// OK
			} else if dotGot {
				intAfterDotGot = true
			} else if negative {
				if integer == 0 && idx != offset+1 {
					err = it.numErrorf(idx, "unexpected zero")
					return
				}
			} else if integer == 0 {
				err = it.numErrorf(idx, "unexpected zero")
				return
			}
			integer *= 10

		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if !floated {
				integer = integer*10 + uint64(b) - '0'
			} else if !exponentGot {
				intAfterDotGot = true
			}

		case 'e', 'E':
			if exponentGot {
				err = it.numErrorf(idx, "unexpected exponent symbol")
				return
			}
			exponentGot = true
			floated = true

		case '+':
			if !exponentGot {
				err = it.numErrorf(idx, "unexpected +")
				return
			}
			// Codes below not needed because this error is caught in outer logic
			// if !floated {
			// 	err = it.numErrorf(idx, "unexpected positive symbol")
			// 	return
			// }

		case '-':
			if !floated {
				if idx != offset {
					err = it.numErrorf(idx, "unexpected negative symbol")
					return
				}
				negative = true
			}

		case '.':
			if idx == offset || floated || exponentGot || dotGot {
				err = it.numErrorf(idx, "unexpected dot symbol")
				return
			}
			dotGot = true
			floated = true
		}
	}

	if edgeFound {
		idx--
	}

	if floated {
		if dotGot && !intAfterDotGot {
			err = it.numErrorf(offset, "integer after dot missing")
			return
		}
		v, err = it.parseFloatResult(p, offset, idx)
	} else {
		if integer > 0 && it[offset] == '0' {
			err = it.numErrorf(offset, "non-zero integer should not start with zero")
			return
		}

		firstB := it[offset]
		if idx-offset == 1 {
			if firstB >= '0' && firstB <= '9' {
				// OK
			} else {
				err = it.numErrorf(offset, "invalid number format")
				return
			}
		}

		if negative {
			v, err = it.parseNegativeIntResult(p, offset, idx, integer)
		} else {
			v, err = it.parsePositiveIntResult(p, offset, idx, integer)
		}
	}

	return v, idx, len(it)-idx == 0, err
}

func (it iter) numErrorf(offset int, f string, a ...any) error {
	a = append([]any{offset}, a...)
	return fmt.Errorf("parsing number at index %d: "+f, a...)

	// debug ONLY below

	// getCaller := func(skip int) string {
	// 	pc, _, _, ok := runtime.Caller(skip + 1)
	// 	if !ok {
	// 		return "<caller N/A>"
	// 	}
	// 	ca := runtime.CallersFrames([]uintptr{pc})
	// 	fr, _ := ca.Next()

	// 	fu := filepath.Ext(fr.Function)
	// 	fu = strings.TrimLeft(fu, ".")
	// 	li := fr.Line

	// 	return fmt.Sprintf("%s(), Line %d", fu, li)
	// }
	// ca := getCaller(1)

	// a = append([]any{ca, string(it), offset}, a...)
	// return fmt.Errorf("%s - parsing number \"%s\" at index %d: "+f, a...)
}

const (
	uintMaxStr    = "18446744073709551615"
	uintMaxDigits = 10000000000000000000
	intMin        = -9223372036854775808
	intMinStr     = "-9223372036854775808"
	intMinAbs     = 9223372036854775808
)

func (it iter) parseFloatResult(p pool, start, end int) (*V, error) {
	f, err := strconv.ParseFloat(unsafe.BtoS(it[start:end]), 64)
	if err != nil {
		return nil, it.numErrorf(start, "%w", err)
	}

	v := new(p, Number)
	v.srcByte = it[start:end]

	v.num.negative = f < 0
	v.num.floated = true
	v.num.i64 = int64(f)
	v.num.u64 = uint64(f)
	v.num.f64 = f

	return v, nil
}

func (it iter) parsePositiveIntResult(p pool, start, end int, integer uint64) (*V, error) {
	le := end - start

	if le > len(uintMaxStr) {
		return nil, it.numErrorf(start, "value too large")
	} else if le == len(uintMaxStr) {
		if integer < uintMaxDigits {
			return nil, it.numErrorf(start, "value too large")
		}
	}

	v := new(p, Number)
	v.srcByte = it[start:end]

	v.num.negative = false
	v.num.floated = false
	v.num.i64 = int64(integer)
	v.num.u64 = uint64(integer)
	v.num.f64 = float64(integer)

	return v, nil
}

func (it iter) parseNegativeIntResult(p pool, start, end int, integer uint64) (*V, error) {
	le := end - start

	if le > len(intMinStr) {
		return nil, it.numErrorf(start, "absolute value too large")
	} else if le == len(intMinStr) {
		if integer > intMinAbs {
			return nil, it.numErrorf(start, "absolute value too large")
		}
	}

	v := new(p, Number)
	v.srcByte = it[start:end]

	v.num.negative = true
	v.num.floated = false

	if integer == intMinAbs {
		v.num.i64 = intMin
	} else {
		v.num.i64 = -int64(integer)
	}

	v.num.u64 = uint64(v.num.i64)
	v.num.f64 = float64(integer)

	return v, nil
}
