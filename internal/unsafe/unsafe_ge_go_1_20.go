//go:build go1.20
// +build go1.20

package unsafe

import (
	"unsafe"
)

// BtoS []byte to string
func BtoS(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	// reference: String method of strings.Builder
	return unsafe.String(&b[0], len(b))
}

// StoB string to []byte
func StoB(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
