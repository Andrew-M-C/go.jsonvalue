package unsafe

import (
	"reflect"
	"unsafe"
)

// BtoS []byte to string
func BtoS(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StoB string to []byte
func StoB(s string) []byte {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	sh.Cap = sh.Len
	return *(*[]byte)(unsafe.Pointer(sh))
}
