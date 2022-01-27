package jsonvalue

import (
	"errors"
	"fmt"
)

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
