// Package buffer implements a marshaling buffer for jsonvalue
package buffer

// Buffer defines a buffer type
type Buffer interface {
	WriteByte(byte) error
	Write(d []byte) (int, error)
	WriteString(s string) (int, error)
	WriteRune(r rune) (int, error)
	Bytes() []byte
}

// NewBuffer returns a buffer
func NewBuffer() Buffer {
	return &buffer{
		buff: make([]byte, 0, 4096),
	}
}
