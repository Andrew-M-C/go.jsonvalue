package jsonvalue

import (
	"unicode/utf8"
)

// Buffer provides interface like *bytes.Buffer
type Buffer interface {
	Bytes() []byte
	WriteByte(byte) (err error)
	WriteRune(rune) (n int, err error)
	Write([]byte) (n int, err error)
	WriteString(string) (n int, err error)
}

const (
	byteBufferBlockSize = 4096 // 4kB
	byteBufferCap       = 4096
)

var _ Buffer = (*bufferImpl)(nil)

type bufferImpl struct {
	sect struct {
		block   []byte
		start   int
		remains int
	}

	allSects [][]byte

	total int
}

func NewBuffer() *bufferImpl {
	b := &bufferImpl{}
	b.resetBlock()
	b.allSects = make([][]byte, 0, byteBufferCap)
	return b
}

func (b *bufferImpl) resetBlock() {
	b.sect.block = make([]byte, byteBufferBlockSize)
	b.sect.remains = byteBufferBlockSize
	b.sect.start = 0
}

func (b *bufferImpl) sectSize() int {
	return len(b.sect.block) - b.sect.remains - b.sect.start
}

func (b *bufferImpl) addBlockToSectsIfNeeded() {
	le := b.sectSize()
	if le == 0 {
		return
	}

	start, end := b.sect.start, b.sect.start+le
	b.allSects = append(b.allSects, b.sect.block[start:end])

	b.sect.start += le
}

// Bytes 输出结果
func (b *bufferImpl) Bytes() []byte {
	res := make([]byte, b.total)
	offset := 0
	for _, sect := range b.allSects {
		copy(res[offset:], sect)
		offset += len(sect)
	}

	size := b.sectSize()
	if size > 0 {
		start := b.sect.start
		end := start + size
		copy(res[offset:], b.sect.block[start:end])
	}
	return res
}

func (b *bufferImpl) WriteByte(by byte) error {
	if b.sect.remains < 1 {
		b.addBlockToSectsIfNeeded()
		b.resetBlock()
	}

	size := b.sectSize()
	b.sect.block[b.sect.start+size] = by

	b.sect.remains--
	b.total++

	return nil
}

func (b *bufferImpl) WriteRune(r rune) (n int, err error) {
	if r < 128 {
		return 1, b.WriteByte(byte(r))
	}

	if b.sect.remains < 3 {
		b.addBlockToSectsIfNeeded()
		b.resetBlock()
	}

	size := b.sectSize()

	cnt := utf8.EncodeRune(b.sect.block[b.sect.start+size:], r)
	b.sect.remains -= cnt
	b.total += cnt
	return cnt, nil
}

func (b *bufferImpl) Write(data []byte) (n int, err error) {
	b.addBlockToSectsIfNeeded()
	b.allSects = append(b.allSects, data)
	b.total += len(data)
	return len(data), nil
}

func (b *bufferImpl) WriteString(s string) (n int, err error) {
	return b.Write(unsafeStoB(s))
}
