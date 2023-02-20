package buffer

import (
	"github.com/Andrew-M-C/go.jsonvalue/internal/unsafe"
)

type buffer struct {
	buff []byte
}

func (b *buffer) WriteByte(c byte) error {
	b.buff = append(b.buff, c)
	return nil
}

func (b *buffer) Write(d []byte) (int, error) {
	b.buff = append(b.buff, d...)
	return len(d), nil
}

func (b *buffer) WriteString(s string) (int, error) {
	d := unsafe.StoB(s)
	return b.Write(d)
}

func (b *buffer) WriteRune(r rune) (int, error) {
	return b.WriteString(string(r))
}

func (b *buffer) Bytes() []byte {
	return b.buff
}
